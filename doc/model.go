package doc

import (
	"fmt"
	"github.com/fatih/structtag"
	"github.com/go-openapi/spec"
	"go/types"
)

type targetType struct {
	name    string
	origObj types.Object
}

func newTargetType(name string, origObj types.Object) *targetType {
	return &targetType{name: name, origObj: origObj}
}

func (t *targetType) isValid() bool {
	return t.origObj != nil && t.origObj.Type() != nil && t.origObj.Type().Underlying() != nil
}

func (t *targetType) isStruct() bool {
	_, ok := t.origObj.Type().Underlying().(*types.Struct)
	return ok
}

func (t *targetType) isNamedType() bool {
	_, ok := t.origObj.Type().(*types.Named)
	return ok
}

func (t *targetType) toStruct() *types.Struct {
	return t.origObj.Type().Underlying().(*types.Struct)
}

func (t *targetType) toType() types.Type {
	return t.origObj.Type()
}

func (t *targetType) toNamedType() *types.Named {
	return t.origObj.Type().(*types.Named)
}

func (t *targetType) toTargetStruct() *targetStruct {
	if !t.isStruct() {
		return nil
	}

	return &targetStruct{
		name:       t.name,
		origType:   t.toType(),
		origStruct: t.toStruct(),
	}
}

type targetStruct struct {
	name       string
	origType   types.Type
	origStruct *types.Struct
}

func newTargetStruct(name string, origType types.Type, origStruct *types.Struct) *targetStruct {
	return &targetStruct{name: name, origType: origType, origStruct: origStruct}
}

func (t *targetStruct) isNamedType() bool {
	_, ok := t.origType.(*types.Named)
	return ok
}

func (t *targetStruct) toNamedType() *types.Named {
	return t.origType.(*types.Named)
}

type targetField struct {
	packageID              string
	structName             string
	fieldTag               string
	fieldName              string
	specField              specField
	additionalProperties   specField
	elem                   types.Type
	isArrayType            bool
	isAdditionalProperties bool
}

func (t *targetField) ID() string {
	return fmt.Sprintf("%s.%s.%s", t.packageID, t.structName, t.fieldName)
}

func (t *targetField) CanonicalFieldName(structTag string) string {
	var fieldName = t.fieldName
	if len(t.fieldTag) > 0 {
		if tags, err := structtag.Parse(t.fieldTag); err == nil {
			for _, tp := range tags.Tags() {
				if tp.Key == structTag {
					fieldName = tp.Name
				}
			}
		}
	}

	return fieldName
}

type specField struct {
	baseType, itemsType, format, ref string
}

func (s specField) isValid() bool {
	return s.format != "" || s.ref != "" || s.baseType != ""
}

func (s specField) toSchemaProp(description string) spec.SchemaProps {
	schemaProps := spec.SchemaProps{
		Format:      s.format,
		Description: description,
	}

	if s.baseType == arrayType {
		var props spec.SchemaProps
		if s.ref != "" {
			props = spec.SchemaProps{Ref: spec.MustCreateRef("#/components/schemas/" + s.ref)}
		} else if s.itemsType != "" {
			props = spec.SchemaProps{Type: []string{s.itemsType}}
		}
		schemaProps.Type = []string{s.baseType}
		schemaProps.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: props}}
	} else {
		if s.ref != "" {
			schemaProps.Ref = spec.MustCreateRef("#/components/schemas/" + s.ref)
		} else {
			schemaProps.Type = []string{s.baseType}
		}
	}

	return schemaProps
}
