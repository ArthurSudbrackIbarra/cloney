package basicoperations

// ListContainsString checks if a given list of strings contains a specific value.
func ListContainsString(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
