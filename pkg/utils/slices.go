package utils

// From https://github.com/life4/genesis

// SliceMap applies F to all elements in slice of T and returns slice of results
func SliceMap[S ~[]T, T any, G any](items S, f func(el T) G) []G {
	result := make([]G, 0, len(items))
	for _, el := range items {
		result = append(result, f(el))
	}
	return result
}

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
