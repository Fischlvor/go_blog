package crypto

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 哈希密码
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPasswordDefault 生成随机密码（用于第三方登录）
func HashPasswordDefault() string {
	rand.Seed(time.Now().UnixNano())
	randomPassword := randomString(32)
	hash, _ := HashPassword(randomPassword)
	return hash
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
