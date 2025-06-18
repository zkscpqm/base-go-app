package common

func MergeHeaders(original map[string]string, new map[string]string) map[string]string {
	if original == nil {
		original = make(map[string]string)
	}
	if new == nil {
		return original
	}
	for k, v := range new {
		original[k] = v
	}
	return original
}
