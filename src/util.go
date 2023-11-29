package src

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSecureRandomString 生成指定长度的安全随机字符串
func GenerateSecureRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 使用 base64 进行编码，确保生成的字符串是可打印的
	randomString := base64.StdEncoding.EncodeToString(randomBytes)

	// 去除可能的末尾填充字符
	return randomString[:length], nil
}
