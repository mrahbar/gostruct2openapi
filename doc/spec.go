package doc

import (
	"github.com/go-openapi/spec"
	"sort"
)

type SpecRegistry map[string]spec.Schema

func (s SpecRegistry) AddSchema(key string, schema spec.Schema) {
	s[key] = schema
}

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
