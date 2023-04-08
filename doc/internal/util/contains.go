package util

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
