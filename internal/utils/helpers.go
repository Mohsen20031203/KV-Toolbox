// internal/utils/helpers.go
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
)

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

// SanitizeString پاک‌سازی رشته
func WriteJsonFile(file *os.File, state interface{}) error {
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(&state); err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}
	return nil
}

func IsValidJSON(data string) bool {
	var js json.RawMessage
	m := json.Unmarshal([]byte(data), &js) == nil
	return m
}

func CheckCondition(rightColumnContent *fyne.Container) bool {
	if len(rightColumnContent.Objects) > 2 {
		return false
	}
	return true
}
