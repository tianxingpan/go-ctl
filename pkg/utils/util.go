// Package utils provides file path dependent methods
package utils

import (
	"bytes"
	goformat "go/format"
	"os"
	"path"
	"text/template"
)

// 文件模板配置结构体
type fileTemplateConfig struct {
	dir             string       // 目录
	subdir          string       // 子目录
	filename        string       // 文件名
	templateName    string       // 模板名称
	category        string       // 文件类型
	templateFile    string       // 模板文件
	builtinTemplate string       // 内置模板
	data            interface{}  // 填充内容
}

// 可能创建文件，如果有创建文件，返回文件句柄
func MaybeCreateFile(dir, subDir, file string) (*os.File, bool, error) {
	err := MkdirIfNotExist(path.Join(dir, subDir))
	if err != nil {
		return nil, false, err
	}
	fPath := path.Join(dir, subDir, file)
	if FileExists(fPath) {
		// TODO:打印日志
		return nil, false, nil
	}

	fp, err := CreateIfNotExist(fPath)
	if err != nil {
		return fp, false, err
	}
	return fp, true, err
}

// 代码格式化
func formatCode(code string) string {
	ret, err := goformat.Source([]byte(code))
	if err != nil {
		return code
	}

	return string(ret)
}

// 生成目标文件
func genTargetFile(c fileTemplateConfig) error {
	fp, created, err := MaybeCreateFile(c.dir, c.subdir, c.filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var text string
	if len(c.category) == 0 || len(c.templateFile) == 0 {
		text = c.builtinTemplate
	} else {
		text, err = LoadTemplate(c.category, c.templateFile, c.builtinTemplate)
		if err != nil {
			return err
		}
	}

	t := template.Must(template.New(c.templateName).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, c.data)
	if err != nil {
		return err
	}

	code := formatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}