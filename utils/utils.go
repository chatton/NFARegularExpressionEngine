package utils

// this method is adapted from the example shown here.
// https://www.geeksforgeeks.org/program-print-substrings-given-string/
func AllSubstrings(str string) []string {
	var all []string
	for i := 0; i < len(str); i++ {
		for j := i + 1; j <= len(str); j++ {
			all = append(all, str[i:j])
		}
	}
	return all
}
