package utils

// Remove removes an element at the specified index from a slice of any type
func Remove[T any](slice []T, idx int) []T {
	if idx < 0 || idx >= len(slice) {
		return slice
	}
	return append(slice[:idx], slice[idx+1:]...)
}
