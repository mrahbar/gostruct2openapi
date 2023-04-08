package util

func RemoveElement[T comparable](arr []T, elem T) []T {
	tmp := arr[:0]
	for _, p := range arr {
		if p != elem {
			tmp = append(tmp, p)
		}
	}
	return tmp
}

func Filter[T any](methods []T, pred func(m T) bool) []T {
	var res []T
	for _, m := range methods {
		if pred(m) {
			res = append(res, m)
		}
	}

	return res
}

func Map[T any, U any](input []T, transform func(m T) U) []U {
	var res []U
	for _, u := range input {
		res = append(res, transform(u))
	}

	return res
}
