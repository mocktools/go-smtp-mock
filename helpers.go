package smtpmock

import "regexp"

// Regex builder
func newRegex(regexPattern string) (*regexp.Regexp, error) {
	return regexp.Compile(regexPattern)
}

// Matches string to regex pattern
func matchRegex(strContext, regexPattern string) bool {
	regex, err := newRegex(regexPattern)
	if err != nil {
		return false
	}

	return regex.MatchString(strContext)
}

// Returns string by regex pattern capture group index
func regexCaptureGroup(str string, regexPattern string, captureGroup int) string {
	regex, _ := newRegex(regexPattern)

	return regex.FindStringSubmatch(str)[captureGroup]
}

// Returns true if the given string is present in slice, otherwise returns false
func isIncluded(slice []string, target string) bool {
	if len(slice) > 0 {
		for _, item := range slice {
			if item == target {
				return true
			}
		}
	}

	return false
}
