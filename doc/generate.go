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
	DocumentStruct(filter *regexp.Regexp, _package ...string) ([]spec.Schema, error)
}

type specField struct {
	baseType, itemsType, format, ref string
}

type openapiGenerator struct {
}

func NewOpenapiGenerator() Generator {
	return &openapiGenerator{}
}

func (o *openapiGenerator) DocumentStruct(filter *regexp.Regexp, _package ...string) ([]spec.Schema, error) {
	pkgs, err := loadPackages(_package...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package %s load failed", _package)
	}

	return o.parse(pkgs, filter), nil
}

func (o *openapiGenerator) parse(pkgs []*packages.Package, filter *regexp.Regexp) (specs []spec.Schema) {
	for _, pkg := range pkgs {
		commentMap := loadCommentMap(pkg, filter)
		scope := pkg.Types.Scope()

		for _, structScopeName := range scope.Names() {
			if !filter.MatchString(structScopeName) {
				continue
			}
			obj := scope.Lookup(structScopeName)
			if obj != nil && obj.Type() != nil && obj.Type().Underlying() != nil {
				if _struct, ok := obj.Type().Underlying().(*types.Struct); ok { //skip if underlying scope is not a struct, e.g. interface
					fmt.Printf("Processing struct: name=%s\n", structScopeName)
					props := spec.SchemaProps{ID: structScopeName, Type: []string{objectType}, Properties: make(spec.SchemaProperties)}
					subSpecs := o.toSpec(&props, _struct, commentMap, structScopeName)
					specs = append(specs, spec.Schema{SchemaProps: props})
					for _, s := range subSpecs {
						specs = append(specs, s)
					}

					_structTyp := obj.Type().(*types.Named)
					for i := 0; i < _structTyp.NumMethods(); i++ {
						scope := _structTyp.Method(i).Scope()
						for _, methodScopeName := range scope.Names() {
							if !filter.MatchString(methodScopeName) {
								continue
							}
							obj := scope.Lookup(methodScopeName) //TODO try to rename ID of methodScopeName to public name using @title
							if obj != nil && obj.Type() != nil && obj.Type().Underlying() != nil {
								if _struct, ok := obj.Type().Underlying().(*types.Struct); ok { //skip if underlying scope is not a struct, e.g. interface
									fmt.Printf("Processing method struct: name=%s\n", methodScopeName)
									props := spec.SchemaProps{ID: methodScopeName, Type: []string{objectType}, Properties: make(spec.SchemaProperties)}
									subSpecs := o.toSpec(&props, _struct, commentMap, methodScopeName)
									specs = append(specs, spec.Schema{SchemaProps: props})
									for _, s := range subSpecs {
										specs = append(specs, s)
									}
								}
							}
						}
					}
				}

			}
		}
	}

	return
}

func (o *openapiGenerator) toSpec(props *spec.SchemaProps, _struct *types.Struct, commentMap map[string]string, structName string) map[string]spec.Schema {
	specs := make(map[string]spec.Schema)

	for i := 0; i < _struct.NumFields(); i++ {
		field := _struct.Field(i)

		if field.Embedded() {
			if _embeddedStruct, ok := field.Type().Underlying().(*types.Struct); ok {
				subSpecs := o.toSpec(props, _embeddedStruct, commentMap, field.Name())
				for k, s := range subSpecs {
					specs[k] = s
				}
			}
		} else if field.Exported() {
			fieldName := field.Name()
			underlying := field.Type().Underlying()

			//early handling of time.Time due to underlying type is actually a struct
			if fieldType := structFieldTypeMap[underlying.String()]; fieldType.format == "RFC3339" {
				o.mapField(props, _struct.Tag(i), commentMap, structName, fieldName, fieldType)
				continue
			}

			switch u := underlying.(type) {
			case *types.Interface:
				// falling back to object type because handling of interface type is not possible
				o.mapField(props, _struct.Tag(i), commentMap, structName, fieldName, specField{baseType: objectType})
			case *types.Struct:
				name := field.Type().(*types.Named).Obj().Name()
				o.mapField(props, _struct.Tag(i), commentMap, structName, fieldName, specField{ref: name})
				subProps := spec.SchemaProps{ID: name, Type: []string{objectType}, Properties: make(spec.SchemaProperties)}
				subSpecs := o.toSpec(&subProps, u, commentMap, name)
				specs[name] = spec.Schema{SchemaProps: subProps}
				for k, s := range subSpecs {
					specs[k] = s
				}
			case *types.Pointer:
				subSpecs := o.handleUnderlyingField(props, commentMap, structName, fieldName, _struct.Tag(i), u.Elem(), false)
				for k, s := range subSpecs {
					specs[k] = s
				}
			case *types.Slice:
				subSpecs := o.handleUnderlyingField(props, commentMap, structName, fieldName, _struct.Tag(i), u.Elem(), true)
				for k, s := range subSpecs {
					specs[k] = s
				}
			default:
				o.mapField(props, _struct.Tag(i), commentMap, structName, fieldName, structFieldTypeMap[underlying.String()])
			}
		}
	}

	return specs
}

func (o *openapiGenerator) handleUnderlyingField(props *spec.SchemaProps, commentMap map[string]string, structName string,
	fieldName string, tag string, elem types.Type, isArrayType bool) map[string]spec.Schema {
	specs := make(map[string]spec.Schema)

	if _underlyingStruct, ok := elem.Underlying().(*types.Pointer); ok {
		return o.handleUnderlyingField(props, commentMap, structName, fieldName, tag, _underlyingStruct.Elem(), isArrayType)
	} else if _underlyingStruct, ok := elem.Underlying().(*types.Struct); ok {
		name := elem.(*types.Named).Obj().Name()
		_specField := specField{ref: name}
		if isArrayType {
			_specField.baseType = arrayType
		}
		o.mapField(props, tag, commentMap, structName, fieldName, _specField)
		subProps := spec.SchemaProps{ID: name, Type: []string{objectType}, Properties: make(spec.SchemaProperties)}
		subSpecs := o.toSpec(&subProps, _underlyingStruct, commentMap, name)
		specs[name] = spec.Schema{SchemaProps: subProps}
		for k, s := range subSpecs {
			specs[k] = s
		}
	} else if _underlyingBasic, ok := elem.Underlying().(*types.Basic); ok {
		var _specField specField
		if isArrayType {
			_specField = specField{baseType: arrayType, itemsType: structFieldTypeMap[_underlyingBasic.String()].baseType}
		} else {
			_specField = structFieldTypeMap[_underlyingBasic.String()]
		}
		o.mapField(props, tag, commentMap, structName, fieldName, _specField)
	} else {
		fmt.Printf("has no basic type but %s", elem.Underlying().String())
		var _specField specField
		if isArrayType {
			_specField = specField{baseType: arrayType, itemsType: objectType}
		} else {
			_specField = specField{baseType: objectType}
		}
		o.mapField(props, tag, commentMap, structName, fieldName, _specField)
	}

	return specs
}

func (o *openapiGenerator) mapField(props *spec.SchemaProps, tag string, commentMap map[string]string, structName, fieldName string, specField specField) {
	description := strings.Replace(commentMap[fmt.Sprintf("%s.%s", structName, fieldName)], "\n", "", -1)

	if len(tag) > 0 {
		tags, err := structtag.Parse(tag)
		if err == nil {
			for _, tp := range tags.Tags() {
				if tp.Key == "json" {
					fieldName = tp.Name
				}
			}
		}
	}

	schemaProps := spec.SchemaProps{
		Format:      specField.format,
		Description: description,
	}

	if specField.baseType == arrayType {
		var props spec.SchemaProps
		if specField.ref != "" {
			props = spec.SchemaProps{Ref: spec.MustCreateRef("#/components/schemas/" + specField.ref)}
		} else if specField.itemsType != "" {
			props = spec.SchemaProps{Type: []string{specField.itemsType}}
		}
		schemaProps.Type = []string{specField.baseType}
		schemaProps.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: props}}
	} else {
		if specField.ref != "" {
			schemaProps.Ref = spec.MustCreateRef("#/components/schemas/" + specField.ref)
		} else {
			schemaProps.Type = []string{specField.baseType}
		}
	}

	props.Properties[fieldName] = spec.Schema{
		SchemaProps: schemaProps,
	}
}
