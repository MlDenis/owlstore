package cache

import (
	"sync"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/pkg"
)

// структура где будут хранится токены в оперативной памяти
type DataJWT struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewDataJWT() *DataJWT {
	return &DataJWT{
		data: map[string]string{},
	}
}

// добавляем токен в нашу структуру
func (userJWT *DataJWT) AddToken(userData *models.UserData) {
	userJWT.mu.Lock()
	defer userJWT.mu.Unlock()
	userJWT.data[userData.Token] = userData.Login

}

// получаем токен(для проверки авторизации)
func (userJWT *DataJWT) GetToken(userData *models.UserData) error {
	login, ok := userJWT.data[userData.Token]
	if !ok {
		return pkg.TokenNotExist
	}
	userData.Login = login
	return nil
}
