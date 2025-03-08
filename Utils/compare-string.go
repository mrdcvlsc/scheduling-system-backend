package Utils

import "strings"

// ignores white spaces and case.
func HasSubString(main string, sub string) bool {
	processed_main := RemoveWhiteSpace(main)
	processed_main = strings.ToLower(processed_main)

	processed_sub := RemoveWhiteSpace(sub)
	processed_sub = strings.ToLower(processed_sub)

	return strings.Contains(processed_main, processed_sub)
}
