package utils

// ArrayIndexOf TODO go2.0 多态
func ArrayIndexOf(arr []string, s string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}
