package finalizer

func ContainsStringHelloSonarone(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
