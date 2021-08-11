package libs

import (
	"regexp"
	"strings"
)

func ValidatePassword(pwd string) bool {
	uppercaseValidator := regexp.MustCompile(`[A-Z]`)
	lowercaseValidator := regexp.MustCompile(`[a-z]`)
	numberValidator := regexp.MustCompile(`[0-9]`)

	return uppercaseValidator.MatchString(pwd) &&
		lowercaseValidator.MatchString(pwd) &&
		numberValidator.MatchString(pwd) &&
		len(pwd) > 8
}

func ValidateUsername(uname string) bool {
	return len(uname) > 8
}

func ValidateUploadFileType(filename string) bool {
	return strings.Contains(filename, ".jpg") ||
		strings.Contains(filename, ".jpeg") ||
		strings.Contains(filename, ".png") ||
		strings.Contains(filename, ".mp4")
}
