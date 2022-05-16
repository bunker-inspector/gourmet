package gourmet

func Collect[T any](seq chan T) []T {
	var result []T
	for v := range seq {
		result = append(result, v)
	}
	return result
}
