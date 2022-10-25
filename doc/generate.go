package doc

import (
	"fmt"
	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
	"strings"
)

type Generator interface {
	DocumentStruct(_package ...string) ([]spec.Schema, error)
}

type specField struct {
	baseType, itemsType, format, ref string
}

type openapiGenerator struct {
	filter     *regexp.Regexp
	structTags []string
}

func NewOpenapiGenerator(filter *regexp.Regexp, tags ...string) Generator {
	return &openapiGenerator{filter: filter, structTags: tags}
}

func (o *openapiGenerator) DocumentStruct(_package ...string) ([]spec.Schema, error) {
	pkgs, err := loadPackages(_package...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package %s load failed", _package)
	}

	registry := o.parse(pkgs)
	return registry.Values(), nil
}

func (o *openapiGenerator) parse(pkgs []*packages.Package) SpecRegistry {
	specs := make(SpecRegistry)

	for _, pkg := range pkgs {
		commentMap := loadCommentMap(pkg)
		scope := pkg.Types.Scope()

		for _, structScopeName := range scope.Names() {
			if !o.filter.MatchString(structScopeName) {
				continue
			}
			specs.Extend(o.processObj(targetType{structScopeName, scope.Lookup(structScopeName)}, commentMap))
		}
	}

	return specs
}

func (o *openapiGenerator) processObj(target targetType, commentMap map[string]string) SpecRegistry {
	if !target.isValid() || !target.isStruct() {
		return nil
	}

	specs := make(SpecRegistry)
	specs.Extend(o.processTarget(target.toTargetStruct(), commentMap))
	if target.isNamedType() {
		specs.Extend(o.processStructMethods(target.toNamedType(), commentMap))
	}

	return specs
}

func (o *openapiGenerator) processTarget(target *targetStruct, commentMap map[string]string) SpecRegistry {
	fmt.Printf("Processing struct: name=%s\n", target.name)

	var props = spec.SchemaProps{ID: target.name, Type: []string{objectType}, Properties: make(spec.SchemaProperties)}
	subSpecs := o.toSpec(&props, target, commentMap)

	specs := make(SpecRegistry)
	specs.AddSchemaProp(target.name, props)
	specs.Extend(subSpecs)

	return specs
}

func (o *openapiGenerator) processStructMethods(_structTyp *types.Named, commentMap map[string]string) SpecRegistry {
	specs := make(SpecRegistry)

	for i := 0; i < _structTyp.NumMethods(); i++ {
		scope := _structTyp.Method(i).Scope()
		for _, methodScopeName := range scope.Names() {
			if !o.filter.MatchString(methodScopeName) {
				continue
			}
			specs = o.processObj(targetType{methodScopeName, scope.Lookup(methodScopeName)}, commentMap)
		}
	}

	return specs
}

func (o *openapiGenerator) toSpec(props *spec.SchemaProps, target *targetStruct, commentMap map[string]string) SpecRegistry {
	specs := make(SpecRegistry)

	for i := 0; i < target.origStruct.NumFields(); i++ {
		field := target.origStruct.Field(i)

		if field.Embedded() {
			if _embeddedStruct, ok := field.Type().Underlying().(*types.Struct); ok {
				subSpecs := o.toSpec(props, newTargetStruct(field.Name(), field.Type(), _embeddedStruct), commentMap)
				specs.Extend(subSpecs)
			}
		} else if field.Exported() {
			fieldName := field.Name()
			underlying := field.Type().Underlying()

			//early handling of time.Time due to underlying type is actually a struct
			if fieldType := structFieldTypeMap[underlying.String()]; fieldType.format == "RFC3339" {
				tf := &targetField{
					structName: target.name,
					fieldTag:   target.origStruct.Tag(i),
					fieldName:  fieldName,
					specField:  fieldType,
				}
				o.mapField(props, tf, commentMap)
				continue
			}

			switch u := underlying.(type) {
			case *types.Interface:
				// falling back to object type because handling of interface type is not possible
				tf := &targetField{
					structName: target.name,
					fieldTag:   target.origStruct.Tag(i),
					fieldName:  fieldName,
					specField:  specField{baseType: objectType},
				}
				o.mapField(props, tf, commentMap)
			case *types.Struct:
				obj := field.Type().(*types.Named).Obj()
				name := obj.Name()
				tf := &targetField{
					structName: target.name,
					fieldTag:   target.origStruct.Tag(i),
					fieldName:  fieldName,
					specField:  specField{ref: name},
				}
				o.mapField(props, tf, commentMap)
				specs.Extend(o.processTarget(newTargetStruct(name, field.Type(), u), commentMap))
			case *types.Pointer:
				tf := &targetField{
					structName: target.name,
					fieldTag:   target.origStruct.Tag(i),
					fieldName:  fieldName,
					elem:       u.Elem(),
				}
				specs.Extend(o.handleUnderlyingField(props, tf, commentMap))
			case *types.Slice:
				tf := &targetField{
					structName:  target.name,
					fieldTag:    target.origStruct.Tag(i),
					fieldName:   fieldName,
					elem:        u.Elem(),
					isArrayType: true,
				}
				specs.Extend(o.handleUnderlyingField(props, tf, commentMap))
			default:
				tf := &targetField{
					structName: target.name,
					fieldTag:   target.origStruct.Tag(i),
					fieldName:  fieldName,
					specField:  structFieldTypeMap[underlying.String()],
				}
				o.mapField(props, tf, commentMap)
			}
		}
	}

	return specs
}

//, structName string, fieldName string, tag string, elem types.Type, isArrayType bool
func (o *openapiGenerator) handleUnderlyingField(props *spec.SchemaProps, target *targetField, commentMap map[string]string) SpecRegistry {
	specs := make(SpecRegistry)

	if _underlyingStruct, ok := target.elem.Underlying().(*types.Pointer); ok {
		target.elem = _underlyingStruct.Elem()
		return o.handleUnderlyingField(props, target, commentMap)
	} else if _underlyingStruct, ok := target.elem.Underlying().(*types.Struct); ok {
		field := target.elem.(*types.Named).Obj()
		name := field.Name()
		_specField := specField{ref: name}
		if target.isArrayType {
			_specField.baseType = arrayType
		}
		target.specField = _specField
		o.mapField(props, target, commentMap)
		specs.Extend(o.processTarget(newTargetStruct(name, field.Type(), _underlyingStruct), commentMap))
	} else if _underlyingBasic, ok := target.elem.Underlying().(*types.Basic); ok {
		var _specField specField
		if target.isArrayType {
			_specField = specField{baseType: arrayType, itemsType: structFieldTypeMap[_underlyingBasic.String()].baseType}
		} else {
			_specField = structFieldTypeMap[_underlyingBasic.String()]
		}
		target.specField = _specField
		o.mapField(props, target, commentMap)
	} else {
		fmt.Printf("has no basic type but %s", target.elem.Underlying().String())
		var _specField specField
		if target.isArrayType {
			_specField = specField{baseType: arrayType, itemsType: objectType}
		} else {
			_specField = specField{baseType: objectType}
		}
		target.specField = _specField
		o.mapField(props, target, commentMap)
	}

	return specs
}

func (o *openapiGenerator) mapField(props *spec.SchemaProps, target *targetField, commentMap map[string]string) {
	description := strings.Replace(commentMap[fmt.Sprintf("%s.%s", target.structName, target.fieldName)], "\n", "", -1)

	var fieldName = target.fieldName
	if len(target.fieldTag) > 0 {
		tags, err := structtag.Parse(target.fieldTag)
		if err == nil {
			for _, tp := range tags.Tags() {
				if Contains(o.structTags, tp.Key) {
					fieldName = tp.Name
				}
			}
		}
	}

	schemaProps := spec.SchemaProps{
		Format:      target.specField.format,
		Description: description,
	}

	if target.specField.baseType == arrayType {
		var props spec.SchemaProps
		if target.specField.ref != "" {
			props = spec.SchemaProps{Ref: spec.MustCreateRef("#/components/schemas/" + target.specField.ref)}
		} else if target.specField.itemsType != "" {
			props = spec.SchemaProps{Type: []string{target.specField.itemsType}}
		}
		schemaProps.Type = []string{target.specField.baseType}
		schemaProps.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: props}}
	} else {
		if target.specField.ref != "" {
			schemaProps.Ref = spec.MustCreateRef("#/components/schemas/" + target.specField.ref)
		} else {
			schemaProps.Type = []string{target.specField.baseType}
		}
	}

	props.Properties[fieldName] = spec.Schema{
		SchemaProps: schemaProps,
	}
}
