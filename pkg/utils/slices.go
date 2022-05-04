package utils

// Map applies F to all elements in slice of T and returns slice of results
// from https://github.com/life4/genesis
func SliceMap[S ~[]T, T any, G any](items S, f func(el T) G) []G {
	result := make([]G, 0, len(items))
	for _, el := range items {
		result = append(result, f(el))
	}
	return result
}
