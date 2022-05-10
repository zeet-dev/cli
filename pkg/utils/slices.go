package utils

// From https://github.com/life4/genesis

// SliceFilter returns slice of T for which F returned true
func SliceFilter[S ~[]T, T any](items S, f func(el T) bool) S {
	result := make([]T, 0, len(items))
	for _, el := range items {
		if f(el) {
			result = append(result, el)
		}
	}
	return result
}
