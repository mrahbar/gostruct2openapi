package doc

import (
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/mrahbar/gostruct2openapi/doc/internal"
	"github.com/mrahbar/gostruct2openapi/doc/internal/util"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
)

const defaultStructTag = "json"

// Generator generated the OpenAPI document for the named packages
type Generator interface {
	DocumentStruct(_package ...string) ([]spec.Schema, error)
}

type openapiGenerator struct {
	filter           *regexp.Regexp
	structTag        string
	commentRegistry  *internal.CommentRegistry
	metadataParser   *internal.MetadataParser
	processedTargets map[string]struct{}
}

// NewOpenapiGenerator returns a new Generator
func NewOpenapiGenerator(filter *regexp.Regexp, structTag string) Generator {
	if len(structTag) == 0 {
		structTag = defaultStructTag
	}

	return &openapiGenerator{
		filter:           filter,
		structTag:        structTag,
		commentRegistry:  internal.NewCommentRegistry(),
		metadataParser:   internal.NewMetadataParser(),
		processedTargets: make(map[string]struct{}),
	}
}

func (o *openapiGenerator) DocumentStruct(_package ...string) ([]spec.Schema, error) {
	pkgs, err := loadPackages(_package...)
	if err != nil {
		return nil, err
	}

	registry := o.parse(pkgs)
	return registry.Values(), nil
}

func (o *openapiGenerator) parse(pkgs []*packages.Package) SpecRegistry {
	specs := make(SpecRegistry)

	for _, pkg := range pkgs {
		// prepare all comments in package
		o.commentRegistry.Load(pkg)

		// for each package iterate all types (structs, (struct) methods, functions, ...)
		scope := pkg.Types.Scope()
		for _, structScopeName := range scope.Names() {
			if o.doFilter(structScopeName) {
				continue
			}
			specs.Extend(o.processObj(internal.NewTargetType(structScopeName, scope.Lookup(structScopeName))))
		}
	}

	return specs
}

func (o *openapiGenerator) doFilter(value string) bool {
	return !o.filter.MatchString(value)
}

func (o *openapiGenerator) processObj(target *internal.TargetType) SpecRegistry {
	if !target.IsValid() || !target.IsStruct() {
		return nil
	}

	specs := make(SpecRegistry)
	specs.Extend(o.processTarget(target.ToTargetStruct()))
	if target.IsNamedType() {
		specs.Extend(o.processStructMethods(target.ToNamedType()))
	}

	return specs
}

func (o *openapiGenerator) processTarget(target *internal.TargetStruct) SpecRegistry {
	specs := make(SpecRegistry)
	if _, exists := o.processedTargets[target.Name()]; exists {
		return specs
	} else {
		o.processedTargets[target.Name()] = struct{}{}
	}

	fmt.Printf("Processing struct: name=%s\n", target.Name())

	if target.IsNamedType() {
		if pkgs, err := loadPackages(target.ToNamedType().Obj().Pkg().Path()); err == nil {
			o.commentRegistry.Load(pkgs...)
		}
	}

	metadata := o.metadataParser.ParseStructDesc(o.commentRegistry.Lookup(target.ID()))
	var props = spec.SchemaProps{ID: metadata.Lookup(internal.TitleAttr, target.Name()), Type: []string{internal.ObjectType.String()}, Description: util.CleanDescription(metadata.Lookup(internal.DescriptionAttr, "")), Properties: make(spec.SchemaProperties)}
	specs.AddSchemaProp(props)
	specs.Extend(o.toSpec(&props, target))

	return specs
}

func (o *openapiGenerator) processStructMethods(_structTyp *types.Named) SpecRegistry {
	specs := make(SpecRegistry)

	for i := 0; i < _structTyp.NumMethods(); i++ {
		scope := _structTyp.Method(i).Scope()
		if scope == nil {
			fmt.Printf("Method %q of struct %s has no associated scope\n", _structTyp.Method(i).Name(), _structTyp.String())
			continue
		}
		for _, methodScopeName := range scope.Names() {
			if o.doFilter(methodScopeName) {
				continue
			}
			specs.Extend(o.processObj(internal.NewTargetType(methodScopeName, scope.Lookup(methodScopeName))))
		}
	}

	return specs
}

func (o *openapiGenerator) toSpec(props *spec.SchemaProps, target *internal.TargetStruct) SpecRegistry {
	specs := make(SpecRegistry)

	for i := 0; i < target.OriginalStruct().NumFields(); i++ {
		field := target.OriginalStruct().Field(i)

		if field.Embedded() {
			if _embeddedStruct, ok := field.Type().Underlying().(*types.Struct); ok {
				subSpecs := o.toSpec(props, internal.NewTargetStruct(field.Name(), field.Type(), _embeddedStruct))
				specs.Extend(subSpecs)
			}
		} else if field.Exported() {
			fieldName := field.Name()
			underlying := field.Type().Underlying()

			tf := internal.NewTargetField(
				field.Pkg().Path(),
				target.Name(),
				target.OriginalStruct().Tag(i),
				fieldName,
			)

			//early handling of time.Time due to underlying type is actually a struct
			if util.IsTimeField(field.Type()) {
				tf.SetSpecField(internal.StructFieldTypeMap["time.Time"])
				o.mapField(props, tf)
				continue
			}

			switch u := underlying.(type) {
			case *types.Map:
				tf.SetSpecField(internal.NewSpecField(internal.ObjectType))
				tf.SetIsAdditionalProperties()
				tf.SetElem(u.Elem())
				specs.Extend(o.handleUnderlyingField(props, tf))
			case *types.Interface, *types.Chan:
				// falling back to object type because handling of type is not possible
				tf.SetSpecField(internal.NewSpecField(internal.ObjectType))
				o.mapField(props, tf)
			case *types.Struct:
				name := field.Type().(*types.Named).Obj().Name()
				tf.SetSpecField(internal.NewStructSpecField(name))
				o.mapField(props, tf)
				specs.Extend(o.processTarget(internal.NewTargetStruct(name, field.Type(), u)))
			case *types.Pointer:
				tf.SetElem(u.Elem())
				specs.Extend(o.handleUnderlyingField(props, tf))
			case *types.Slice:
				tf.SetElem(u.Elem())
				tf.SetIsArrayType()
				specs.Extend(o.handleUnderlyingField(props, tf))
			default:
				tf.SetSpecField(internal.StructFieldTypeMap[underlying.String()])
				o.mapField(props, tf)
			}
		}
	}

	return specs
}

func (o *openapiGenerator) handleUnderlyingField(props *spec.SchemaProps, target *internal.TargetField) SpecRegistry {
	specs := make(SpecRegistry)

	switch u := target.UnderlyingElem().(type) {
	case *types.Pointer:
		target.SetElem(u.Elem())
		return o.handleUnderlyingField(props, target)
	case *types.Struct:
		field := target.Elem().(*types.Named).Obj()
		name := field.Name()
		var sf *internal.SpecField
		if target.IsArrayType() {
			sf = internal.NewArraySpecField(internal.StructType)
			sf.SetRef(name)
		} else {
			sf = internal.NewStructSpecField(name)
		}
		if target.IsAdditionalProperties() {
			target.SetAdditionalProperties(sf)
		} else {
			target.SetSpecField(sf)
		}
		o.mapField(props, target)
		specs.Extend(o.processTarget(internal.NewTargetStruct(name, field.Type(), u)))
	case *types.Basic:
		var sf *internal.SpecField
		if target.IsArrayType() {
			sf = internal.NewArraySpecField(internal.StructFieldTypeMap[u.String()].BaseType())
		} else {
			sf = internal.StructFieldTypeMap[u.String()]
		}
		if target.IsAdditionalProperties() {
			target.SetAdditionalProperties(sf)
		} else {
			target.SetSpecField(sf)
		}
		o.mapField(props, target)
	default:
		fmt.Printf("%s has no well-known basic type. Got %s\n", target.ID(), target.UnderlyingElem().String())
		var sf *internal.SpecField
		if target.IsArrayType() {
			sf = internal.NewArraySpecField(internal.ObjectType)
		} else {
			sf = internal.NewSpecField(internal.ObjectType)
		}
		if target.IsAdditionalProperties() {
			target.SetAdditionalProperties(sf)
		} else {
			target.SetSpecField(sf)
		}
		o.mapField(props, target)
	}

	return specs
}

func (o *openapiGenerator) mapField(props *spec.SchemaProps, target *internal.TargetField) {
	schema := spec.Schema{
		SchemaProps: target.SpecField().ToSchemaProp(util.CleanDescription(o.commentRegistry.Lookup(target.ID()))),
	}
	if target.HasAdditionalProperties() {
		schema.AdditionalProperties = &spec.SchemaOrBool{
			Schema: &spec.Schema{
				SchemaProps: target.AdditionalProperties().ToSchemaProp(""),
			},
		}
	}
	props.Properties[target.CanonicalFieldName(o.structTag)] = schema
}
