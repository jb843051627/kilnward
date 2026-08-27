package validation

func UniqueStrings(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value == "" {
			return false
		}
		if _, ok := seen[value]; ok {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func Consecutive(values []int) bool {
	for i, value := range values {
		if value != i {
			return false
		}
	}
	return true
}
