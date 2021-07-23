// Package utils provides file path dependent methods
package utils

import (
	"bytes"
	"strings"
	"unicode"
)

// 自定义的String结构体
type String struct {
	source string
}

// 判断字符串是否为空
func (s String) IsEmptyOrSpace() bool {
	if len(s.source) == 0 {
		return true
	}
	if strings.TrimSpace(s.source) == "" {
		return true
	}
	return false
}

// 转小写
func (s String) Lower() string {
	return strings.ToLower(s.source)
}

// 全部替换
func (s String) ReplaceAll(old, new string) string {
	return strings.ReplaceAll(s.source, old, new)
}

// 源字符串
func (s String) Source() string {
	return s.source
}

// 返回单词首字母大写的拷贝字符串
func (s String) Title() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	return strings.Title(s.source)
}

// 蛇形转驼峰格式
func (s String) ToCamel() string {
	list := s.splitBy(func(r rune) bool {
		return r == '_'
	}, true)
	var target []string
	for _, item := range list {
		target = append(target, StringFrom(item).Title())
	}
	return strings.Join(target, "")
}

// 驼峰转蛇形格式
func (s String) ToSnake() string {
	list := s.splitBy(unicode.IsUpper, false)
	var target []string
	for _, item := range list {
		target = append(target, StringFrom(item).Lower())
	}
	return strings.Join(target, "_")
}

// return original string if rune is not letter at index 0
func (s String) Untitled() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	r := rune(s.source[0])
	if !unicode.IsUpper(r) && !unicode.IsLower(r) {
		return s.source
	}
	return string(unicode.ToLower(r)) + s.source[1:]
}

// 转大写
func (s String) Upper() string {
	return strings.ToUpper(s.source)
}

// it will not ignore spaces
func (s String) splitBy(fn func(r rune) bool, remove bool) []string {
	if s.IsEmptyOrSpace() {
		return nil
	}
	var list []string
	buffer := new(bytes.Buffer)
	for _, r := range s.source {
		if fn(r) {
			if buffer.Len() != 0 {
				list = append(list, buffer.String())
				buffer.Reset()
			}
			if !remove {
				buffer.WriteRune(r)
			}
			continue
		}
		buffer.WriteRune(r)
	}
	if buffer.Len() != 0 {
		list = append(list, buffer.String())
	}
	return list
}

// StringFrom
func StringFrom(data string) String {
	return String{source: data}
}
