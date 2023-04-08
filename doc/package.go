package doc

import (
	"fmt"
	"go/token"
	"golang.org/x/tools/go/packages"
)

// loadPackages loads and returns the named Go packages
func loadPackages(_package ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{Fset: token.NewFileSet(), Mode: packages.NeedTypes | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, _package...)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package %s Load failed", _package)
	}
	return pkgs, nil
}
