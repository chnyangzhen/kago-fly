package helper

import (
	"github.com/google/uuid"
	"strings"
	"unicode"
)

// ContainsString 判断字符串是否在字符串数组中
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func Uuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

/*
字符串两端自动拼接单引号
*/
func QuoteString(s string) string {
	return "'" + s + "'"
}

// 判断字符串s是否包含中文字符
func ContainsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
