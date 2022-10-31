package doc

import (
	"go/types"
	"strings"
)

func IndexOf[T comparable](arr []T, elem T) int {
	for i := range arr {
		if arr[i] == elem {
			return i
		}
	}
	return -1
}

func Contains[T comparable](arr []T, elem T) bool {
	return arr != nil && IndexOf(arr, elem) != -1
}

func Deduplicate[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
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
