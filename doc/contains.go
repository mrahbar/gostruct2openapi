package doc

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
