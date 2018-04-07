package utils

// this method is adapted from the example shown here.
// https://www.geeksforgeeks.org/program-print-substrings-given-string/
func AllSubstrings(str string) []string {
	addedEmptyString := false
	var all []string
	for i := 0; i < len(str); i++ {
		for j := i; j <= len(str); j++ {
			s := str[i:j] // only want to add one empty string so check and don't add second.
			if s == "" && !addedEmptyString {
				addedEmptyString = true
			} else {
				all = append(all, s)
			}
		}
	}
	return all
}
