package util

import (
	"go/types"
	"strings"
)

func CleanDescription(desc string) string {
	return strings.Replace(desc, "\n", "", -1)
}

func IsTimeField(field types.Type) bool {
	switch u := field.(type) {
	case *types.Named:
		return u.Obj().Name() == "Time" && u.Obj().Pkg().Name() == "time"
	case *types.Pointer:
		return IsTimeField(u.Elem())
	}

	return false
}
