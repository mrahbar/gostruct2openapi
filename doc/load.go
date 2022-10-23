package doc

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"regexp"
)

func loadPackages(_package ...string) ([]*packages.Package, error) {
	mode := packages.NeedTypes | packages.NeedSyntax
	cfg := &packages.Config{Fset: token.NewFileSet(), Mode: mode}
	pkgs, err := packages.Load(cfg, _package...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package %s load failed", _package)
	}
	return pkgs, nil
}

func loadCommentMap(pkg *packages.Package, filter *regexp.Regexp) map[string]string {
	commentMap := make(map[string]string)
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
									commentMap[fmt.Sprintf("%s.%s", structName, name.Name)] = f.Doc.Text()
								}
							}
						}
					}
				}
			}
		}
	}

	return commentMap
}

func lookupCommentMap(comments []*ast.CommentGroup, scope *types.Scope) (res []string) {
	if scope == nil {
		return
	}

	for _, comment := range comments {
		i := int(comment.Pos())
		i2 := int(scope.Pos())
		i3 := int(comment.End())
		b := i2 <= i
		b2 := i2 <= i3
		if b && b2 {
			res = append(res, comment.Text())
		}
	}

	return
}
