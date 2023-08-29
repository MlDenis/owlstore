package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

const SECRET_KEY = "secret"

// хэшируем пароль для последующей ей записи в бд
func HashPassword(data string) (hash string) {
	key := []byte(SECRET_KEY)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return hash
}
