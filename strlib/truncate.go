package strlib

// Truncate tries to append the "suffix" to the "str". If the length of the appended string exceeds "n",
// the function truncates the "str" to make sure the "suffix" is appended
func Truncate(str, suffix string, n int) string {
	s := str + suffix
	if len(s) <= n {
		return s
	}
	return s[:len(str)-(len(s)-n)] + suffix
}
