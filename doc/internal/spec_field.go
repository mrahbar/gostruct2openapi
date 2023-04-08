package internal

import "github.com/go-openapi/spec"

type SpecField struct {
	baseType, itemsType SpecType
	format, ref         string
}

func NewSpecFieldWithFormat(baseType SpecType, format string) *SpecField {
	return &SpecField{baseType: baseType, format: format}
}

func NewSpecField(baseType SpecType) *SpecField {
	return &SpecField{baseType: baseType}
}

func NewArraySpecField(itemsType SpecType) *SpecField {
	return &SpecField{baseType: ArrayType, itemsType: itemsType}
}

func NewStructSpecField(ref string) *SpecField {
	return &SpecField{baseType: StructType, ref: ref}
}

func (s *SpecField) BaseType() SpecType {
	return s.baseType
}

func (s *SpecField) SetItemsType(itemsType SpecType) {
	s.itemsType = itemsType
}

func (s *SpecField) SetFormat(format string) {
	s.format = format
}

func (s *SpecField) SetRef(ref string) {
	s.ref = ref
}

func (s *SpecField) IsValid() bool {
	return s.format != "" || s.ref != "" || s.baseType != ""
}

func (s *SpecField) ToSchemaProp(description string) spec.SchemaProps {
	schemaProps := spec.SchemaProps{
		Format:      s.format,
		Description: description,
	}

	if s.baseType == ArrayType {
		var props spec.SchemaProps
		if s.ref != "" {
			props = spec.SchemaProps{Ref: spec.MustCreateRef("#/components/schemas/" + s.ref)}
		} else if s.itemsType != "" {
			props = spec.SchemaProps{Type: []string{s.itemsType.String()}}
		}
		schemaProps.Type = []string{s.baseType.String()}
		schemaProps.Items = &spec.SchemaOrArray{Schema: &spec.Schema{SchemaProps: props}}
	} else {
		if s.ref != "" {
			schemaProps.Ref = spec.MustCreateRef("#/components/schemas/" + s.ref)
			schemaProps.Description = "" //Property 'description' is not allowed for $ref
		} else {
			schemaProps.Type = []string{s.baseType.String()}
		}
	}

	return schemaProps
}
