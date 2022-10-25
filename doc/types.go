package doc

import (
	"github.com/go-openapi/spec"
	"sort"
)

const (
	arrayType   = "array"
	objectType  = "object"
	booleanType = "boolean"
	integerType = "integer"
	numberType  = "number"
	stringType  = "string"

	timeFormat = "RFC3339"
)

var structFieldTypeMap = map[string]specField{
	"string":  {baseType: stringType},
	"int":     {baseType: integerType},
	"float32": {baseType: numberType},
	"float64": {baseType: numberType},
	"bool":    {baseType: booleanType},
	//string representation of time.Time
	"struct{wall uint64; ext int64; loc *time.Location}": {baseType: stringType, format: timeFormat},
}

type SpecRegistry map[string]spec.Schema

func (s SpecRegistry) AddSchemaProp(key string, props spec.SchemaProps) {
	s[key] = spec.Schema{SchemaProps: props}
}

func (s SpecRegistry) Extend(r SpecRegistry) {
	for k, v := range r {
		s[k] = v
	}
}

func (s SpecRegistry) Values() (specs []spec.Schema) {
	for _, v := range s {
		specs = append(specs, v)
	}

	sort.Slice(specs, func(i, j int) bool {
		return specs[i].ID < specs[j].ID
	})
	return
}
