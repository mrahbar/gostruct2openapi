package doc

import (
	"go/types"
	"strings"
)

func IndexOf(arr []string, elem string) int {
	for i := range arr {
		if arr[i] == elem {
			return i
		}
	}
	return -1
}

func Contains(arr []string, elem string) bool {
	return arr != nil && IndexOf(arr, elem) != -1
}

func cleanDescription(desc string) string {
	return strings.Replace(desc, "\n", "", -1)
}

func isTimeField(field types.Type) bool {
	switch u := field.(type) {
	case *types.Named:
		return u.Obj().Name() == "Time" && u.Obj().Pkg().Name() == "time"
	case *types.Pointer:
		return isTimeField(u.Elem())
	}

	return false
}
