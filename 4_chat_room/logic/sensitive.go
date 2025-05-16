package logic

import (
	"4_chat_room/global"
	"strings"
)

// logic/sensitive.go
func FilterSensitive(content string) string {
	for _, word := range global.SensitiveWords {
		content = strings.ReplaceAll(content, word, "**")
	}

	return content
}
