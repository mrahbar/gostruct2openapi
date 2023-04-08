package internal

import (
	"fmt"
	"github.com/mrahbar/gostruct2openapi/doc/internal/util"
	"go/ast"
	"go/doc"
	"golang.org/x/tools/go/packages"
)

type CommentRegistry struct {
	loadedPackages []string
	registry       map[string]string
}

func NewCommentRegistry() *CommentRegistry {
	return &CommentRegistry{registry: make(map[string]string)}
}

// Load loads struct as well as struct field comments and builds comment registry for given packages.
func (c *CommentRegistry) Load(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		if util.Contains(c.loadedPackages, pkg.ID) {
			continue
		}
		c.loadedPackages = append(c.loadedPackages, pkg.ID)
		c.loadStructComments(pkg)
		c.loadStructFieldComments(pkg)
	}
}

func (c *CommentRegistry) loadStructComments(pkg *packages.Package) {
	//transform package.Package to ast.Package
	//note that only the necessary fields are set used by go/doc
	a := &ast.Package{Name: pkg.ID, Files: make(map[string]*ast.File)}
	for k, s := range pkg.Syntax {
		a.Files[fmt.Sprintf("%s_%d", s.Name.String(), k)] = s
	}

	p := doc.New(a, ".", doc.AllDecls)
	for _, t := range p.Types {
		if len(t.Doc) > 0 {
			c.registry[fmt.Sprintf("%s.%s", pkg.ID, t.Name)] = t.Doc
		}
	}
}

func (c *CommentRegistry) loadStructFieldComments(pkg *packages.Package) {
	for _, syntax := range pkg.Syntax {
		for structName, object := range syntax.Scope.Objects {
			switch t := object.Decl.(type) {
			case *ast.TypeSpec:
				switch _struct := t.Type.(type) {
				case *ast.StructType:
					for _, field := range _struct.Fields.List {
						switch field.Type.(type) {
						case *ast.Ident,
							*ast.ArrayType,
							*ast.MapType,
							*ast.ChanType,
							*ast.InterfaceType,
							*ast.StructType,
							// Pointer types are represented via StarExpr nodes.
							*ast.StarExpr:
							for _, name := range field.Names {
								if f, ok := name.Obj.Decl.(*ast.Field); ok && len(f.Doc.Text()) > 0 {
									tf := &TargetField{fieldName: name.Name, structName: structName, packageID: pkg.ID}
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

func (c *CommentRegistry) Lookup(key string) string {
	return c.registry[key]
}
