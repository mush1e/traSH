package utils

// Remove removes an element at the specified index from a slice of any type
func Remove[T any](slice []T, idx int) []T {
	if idx < 0 || idx >= len(slice) {
		return slice
	}
	return append(slice[:idx], slice[idx+1:]...)
}

// maps colors to color codes
func Colorize(text, color string) string {
	colors := map[string]string{
		"black":   "\033[30m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
		"reset":   "\033[0m",
	}
	c, ok := colors[color]
	if !ok {
		c = colors["reset"]
	}
	return c + text + colors["reset"]
}
