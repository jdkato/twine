package tag

// min returns the minimum of `a` and `b`.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// stringInSlice determines if `slice` contains the string `a`.
func stringInSlice(a string, slice []string) bool {
	for _, b := range slice {
		if a == b {
			return true
		}
	}
	return false
}
