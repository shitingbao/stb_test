package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ImageToBase64(imagePath string) (string, error) {
	// 读取图片内容
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	// 获取 MIME 类型前缀（根据扩展名判断）
	ext := strings.ToLower(filepath.Ext(imagePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	default:
		mimeType = "application/octet-stream"
	}

	// 编码为 base64 并拼接成完整的 data URI
	base64Str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str), nil
}

func Base64ToImage(base64Str string) {
	// 示例：你的 Base64 字符串（可以来自前端）
	// base64Str := `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...`

	// 去除前缀（如果有）

	// 去除 data URI 前缀
	if idx := strings.Index(base64Str, "base64,"); idx != -1 {
		base64Str = base64Str[idx+7:]
	}

	// 清理所有空白字符（空格、换行）
	re := regexp.MustCompile(`\s`)
	base64Str = re.ReplaceAllString(base64Str, "")

	// 解码
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Println("解码失败:", err)
		return
	}

	// 写入文件
	err = os.WriteFile("output.png", data, 0644)
	if err != nil {
		fmt.Println("保存失败:", err)
		return
	}

	fmt.Println("保存成功：output.png")
}
