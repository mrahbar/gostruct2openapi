package internal

import (
	"fmt"
	"github.com/fatih/structtag"
	"go/types"
)

type TargetField struct {
	packageID              string
	structName             string
	fieldTag               string
	fieldName              string
	specField              *SpecField
	additionalProperties   *SpecField
	elem                   types.Type
	isArrayType            bool
	isAdditionalProperties bool
}

func NewTargetField(packageID string, structName string, fieldTag string, fieldName string) *TargetField {
	return &TargetField{packageID: packageID, structName: structName, fieldTag: fieldTag, fieldName: fieldName}
}

func (t *TargetField) SpecField() *SpecField {
	return t.specField
}

func (t *TargetField) AdditionalProperties() *SpecField {
	return t.additionalProperties
}

func (t *TargetField) Elem() types.Type {
	return t.elem
}

func (t *TargetField) UnderlyingElem() types.Type {
	return t.elem.Underlying()
}

func (t *TargetField) SetElem(elem types.Type) {
	t.elem = elem
}

func (t *TargetField) IsArrayType() bool {
	return t.isArrayType
}

func (t *TargetField) SetIsArrayType() {
	t.isArrayType = true
}

func (t *TargetField) SetAdditionalProperties(additionalProperties *SpecField) {
	t.additionalProperties = additionalProperties
}

func (t *TargetField) HasAdditionalProperties() bool {
	return t.isAdditionalProperties && t.additionalProperties != nil && t.additionalProperties.IsValid()
}

func (t *TargetField) IsAdditionalProperties() bool {
	return t.isAdditionalProperties
}

func (t *TargetField) SetIsAdditionalProperties() {
	t.isAdditionalProperties = true
}

func (t *TargetField) SetSpecField(specField *SpecField) {
	t.specField = specField
}

func (t *TargetField) ID() string {
	return fmt.Sprintf("%s.%s.%s", t.packageID, t.structName, t.fieldName)
}

func (t *TargetField) CanonicalFieldName(structTag string) string {
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
