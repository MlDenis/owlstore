package balance

import (
	"github.com/MlDenis/internal/gofermart/auth/cache"
	"github.com/MlDenis/internal/gofermart/storage"
)

// структура для наших хэндлеров, далее надо будет добавить возмонжо логер и тд
type HandlerBalanceDB struct {
	StorageBalance storage.InterfaceBalance

	DataJWT *cache.DataJWT
}

func HandlerBalance(balance storage.InterfaceBalance, DataJWT *cache.DataJWT) *HandlerBalanceDB {
	return &HandlerBalanceDB{
		StorageBalance: balance,
		DataJWT:        DataJWT,
	}
}
