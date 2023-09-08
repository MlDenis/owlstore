package users

import (
	"github.com/MlDenis/internal/gofermart/auth/cache"
	"github.com/MlDenis/internal/gofermart/storage"
)

// структура для наших хэндлеров, далее надо будет добавить возмонжо логер и тд
type HandlerUserDB struct {
	StorageUsers storage.InterfaceUser
	DataJWT      *cache.DataJWT
}

func HandlerUsers(users storage.InterfaceUser, DataJWT *cache.DataJWT) *HandlerUserDB {
	return &HandlerUserDB{
		StorageUsers: users,
		DataJWT:      DataJWT,
	}
}
