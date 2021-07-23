// Package utils provides file path dependent methods
package utils

import (
	"strings"
)

const (
	NamingLower string = "lower"
	NamingCamel string = "camel"
	NamingSnake string = "snake"
)

// 命名风格是否有效
func IsNamingValid(name string) (string, bool) {
	if len(name) == 0 {
		name = NamingLower
	}
	switch name {
	case NamingLower, NamingCamel, NamingSnake:
		return name, true
	default:
		return "", false
	}
}

// 格式化文件名
func FormatFilename(filename string, style string) string {
	switch style {
	case NamingCamel:
		return StringFrom(filename).ToCamel()
	case NamingSnake:
		return StringFrom(filename).ToSnake()
	default:
		return strings.ToLower(StringFrom(filename).ToCamel())
	}
}