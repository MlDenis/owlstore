package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// токен будет жить 12 часов
const TOKEN_EXP = time.Hour * 12

// структура для нашего токена, можно добавить условия при необходимости
type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

// создаем токен
func CreateJwtToken(loginUser string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		// собственное утверждение
		Username: loginUser,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
