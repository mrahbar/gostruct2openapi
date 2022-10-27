package doc

import (
	"fmt"
	"github.com/go-openapi/spec"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
)

type Generator interface {
	DocumentStruct(_package ...string) ([]spec.Schema, error)
}

type openapiGenerator struct {
	filter          *regexp.Regexp
	structTag       string
	commentRegistry *CommentRegistry
}

func NewOpenapiGenerator(filter *regexp.Regexp, tags string) Generator {
	if len(tags) == 0 {
		tags = "json"
	}

	return &openapiGenerator{filter: filter, structTag: tags, commentRegistry: newCommentRegistry()}
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
		o.commentRegistry.load(pkg)
		scope := pkg.Types.Scope()

		for _, structScopeName := range scope.Names() {
			if !o.filter.MatchString(structScopeName) {
				continue
			}
			specs.Extend(o.processObj(newTargetType(structScopeName, scope.Lookup(structScopeName))))
		}
	}

	return specs
}

func (o *openapiGenerator) processObj(target *targetType) SpecRegistry {
	if !target.isValid() || !target.isStruct() {
		return nil
	}

	specs := make(SpecRegistry)
	specs.Extend(o.processTarget(target.toTargetStruct()))
	if target.isNamedType() {
		specs.Extend(o.processStructMethods(target.toNamedType()))
	}

	return specs
}

func (o *openapiGenerator) processTarget(target *targetStruct) SpecRegistry {
	fmt.Printf("Processing struct: name=%s\n", target.name)

	if target.isNamedType() {
		if pkgs, err := loadPackages(target.toNamedType().Obj().Pkg().Path()); err == nil {
			o.commentRegistry.load(pkgs...)
		}
	}
	var props = spec.SchemaProps{ID: target.name, Type: []string{objectType}, Description: o.commentRegistry.lookup(target.ID()), Properties: make(spec.SchemaProperties)}
	specs := make(SpecRegistry)
	specs.AddSchemaProp(target.name, props)
	specs.Extend(o.toSpec(&props, target))

	return specs
}

func (o *openapiGenerator) processStructMethods(_structTyp *types.Named) SpecRegistry {
	specs := make(SpecRegistry)

	for i := 0; i < _structTyp.NumMethods(); i++ {
		scope := _structTyp.Method(i).Scope()
		for _, methodScopeName := range scope.Names() {
			if !o.filter.MatchString(methodScopeName) {
				continue
			}
			specs.Extend(o.processObj(newTargetType(methodScopeName, scope.Lookup(methodScopeName))))
		}
	}

	return specs
}

func (o *openapiGenerator) toSpec(props *spec.SchemaProps, target *targetStruct) SpecRegistry {
	specs := make(SpecRegistry)

	for i := 0; i < target.origStruct.NumFields(); i++ {
		field := target.origStruct.Field(i)

		if field.Embedded() {
			if _embeddedStruct, ok := field.Type().Underlying().(*types.Struct); ok {
				subSpecs := o.toSpec(props, newTargetStruct(field.Name(), field.Type(), _embeddedStruct))
				specs.Extend(subSpecs)
			}
		} else if field.Exported() {
			fieldName := field.Name()
			underlying := field.Type().Underlying()

			tf := &targetField{
				packageID:  field.Pkg().Path(),
				structName: target.name,
				fieldTag:   target.origStruct.Tag(i),
				fieldName:  fieldName,
			}

			//early handling of time.Time due to underlying type is actually a struct
			if o.isTimeField(field.Type()) {
				tf.specField = structFieldTypeMap["time.Time"]
				o.mapField(props, tf)
				continue
			}

			switch u := underlying.(type) {
			case *types.Map:
				tf.specField = specField{baseType: objectType}
				tf.isAdditionalProperties = true
				tf.elem = u.Elem()
				specs.Extend(o.handleUnderlyingField(props, tf))
			case *types.Interface, *types.Chan:
				// falling back to object type because handling of type is not possible
				tf.specField = specField{baseType: objectType}
				o.mapField(props, tf)
			case *types.Struct:
				name := field.Type().(*types.Named).Obj().Name()
				tf.specField = specField{ref: name}
				o.mapField(props, tf)
				specs.Extend(o.processTarget(newTargetStruct(name, field.Type(), u)))
			case *types.Pointer:
				tf.elem = u.Elem()
				specs.Extend(o.handleUnderlyingField(props, tf))
			case *types.Slice:
				tf.elem = u.Elem()
				tf.isArrayType = true
				specs.Extend(o.handleUnderlyingField(props, tf))
			default:
				tf.specField = structFieldTypeMap[underlying.String()]
				o.mapField(props, tf)
			}
		}
	}

	return specs
}

func (o *openapiGenerator) handleUnderlyingField(props *spec.SchemaProps, target *targetField) SpecRegistry {
	specs := make(SpecRegistry)

	switch u := target.elem.Underlying().(type) {
	case *types.Pointer:
		target.elem = u.Elem()
		return o.handleUnderlyingField(props, target)
	case *types.Struct:
		field := target.elem.(*types.Named).Obj()
		name := field.Name()
		sf := specField{ref: name}
		if target.isArrayType {
			sf.baseType = arrayType
		}
		if target.isAdditionalProperties {
			target.additionalProperties = sf
		} else {
			target.specField = sf
		}
		o.mapField(props, target)
		specs.Extend(o.processTarget(newTargetStruct(name, field.Type(), u)))
	case *types.Basic:
		var sf specField
		if target.isArrayType {
			sf = specField{baseType: arrayType, itemsType: structFieldTypeMap[u.String()].baseType}
		} else {
			sf = structFieldTypeMap[u.String()]
		}
		if target.isAdditionalProperties {
			target.additionalProperties = sf
		} else {
			target.specField = sf
		}
		o.mapField(props, target)
	default:
		fmt.Printf("has no basic type but %s", target.elem.Underlying().String())
		var sf specField
		if target.isArrayType {
			sf = specField{baseType: arrayType, itemsType: objectType}
		} else {
			sf = specField{baseType: objectType}
		}
		if target.isAdditionalProperties {
			target.additionalProperties = sf
		} else {
			target.specField = sf
		}
		o.mapField(props, target)
	}

	return specs
}

func (o *openapiGenerator) mapField(props *spec.SchemaProps, target *targetField) {
	id := target.ID()
	lookup := o.commentRegistry.lookup(id)
	schema := spec.Schema{
		SchemaProps: target.specField.toSchemaProp(lookup),
	}
	if target.additionalProperties.isValid() {
		schema.AdditionalProperties = &spec.SchemaOrBool{
			Schema: &spec.Schema{
				SchemaProps: target.additionalProperties.toSchemaProp(""),
			},
		}
	}
	props.Properties[target.CanonicalFieldName(o.structTag)] = schema
}

func (o *openapiGenerator) isTimeField(field types.Type) bool {
	switch u := field.(type) {
	case *types.Named:
		return u.Obj().Name() == "Time" && u.Obj().Pkg().Name() == "time"
	case *types.Pointer:
		return o.isTimeField(u.Elem())
	}

	return false
}
