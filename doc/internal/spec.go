package internal

import (
	"github.com/go-openapi/spec"
	"sort"
)

// SpecRegistry holds spec.Schema registered by a given key
type SpecRegistry map[string]spec.Schema

// AddSchema register a spec.Schema registered by a given key
func (s SpecRegistry) AddSchema(key string, schema spec.Schema) {
	s[key] = schema
}

// AddSchemaProp is a convenience methods to call AddSchema for a spec.SchemaProps
// for which a spec.Schema is created and key is derived from SchemaProps ID
func (s SpecRegistry) AddSchemaProp(props spec.SchemaProps) {
	s.AddSchema(props.ID, spec.Schema{SchemaProps: props})
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
