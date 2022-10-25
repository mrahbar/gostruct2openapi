package doc

import (
	"github.com/go-openapi/spec"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"sort"
	"strings"
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
	"string":    {baseType: stringType},
	"int":       {baseType: integerType},
	"float32":   {baseType: numberType},
	"float64":   {baseType: numberType},
	"bool":      {baseType: booleanType},
	"time.Time": {baseType: stringType, format: timeFormat},
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

type CommentRegistry struct {
	loadedPackages []string
	registry       map[string]string
}

func newCommentRegistry() *CommentRegistry {
	return &CommentRegistry{registry: make(map[string]string)}
}

func (c *CommentRegistry) load(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		if Contains(c.loadedPackages, pkg.ID) {
			continue
		}
		c.loadedPackages = append(c.loadedPackages, pkg.ID)
		for _, syntax := range pkg.Syntax {
			for structName, object := range syntax.Scope.Objects {
				switch t := object.Decl.(type) {
				case *ast.TypeSpec:
					switch _struct := t.Type.(type) {
					case *ast.StructType:
						for _, field := range _struct.Fields.List {
							switch field.Type.(type) {
							case *ast.Ident, *ast.ArrayType:
								for _, name := range field.Names {
									if f, ok := name.Obj.Decl.(*ast.Field); ok && len(f.Doc.Text()) > 0 {
										tf := &targetField{fieldName: name.Name, structName: structName, packageID: pkg.ID}
										c.registry[tf.ID()] = f.Doc.Text()
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func (c *CommentRegistry) lookup(key string) string {
	return strings.Replace(c.registry[key], "\n", "", -1)
}
