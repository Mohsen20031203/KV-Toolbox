// internal/utils/helpers.go
package utils

import "strings"

// TruncateString کوتاه کردن رشته به طول مشخص
func TruncateString(input string, length int) string {
	if len(input) > length {
		return input[:length] + "..."
	}
	return input
}

// SanitizeString پاک‌سازی رشته
func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}
