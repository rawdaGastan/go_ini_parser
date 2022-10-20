package ini

func contains(array []string, element string) bool {
	for _, x := range array {
		if x == element {
			return true
		}
	}
	return false
}
