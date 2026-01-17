package utils

import (
	"crypto/sha256"
	"fmt"
)

// HashToken возвращает SHA-256 хеш токена в виде HEX-строки
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash) // 64-символьная hex-строка
}
