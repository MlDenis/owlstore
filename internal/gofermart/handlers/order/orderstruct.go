package order

import (
	"github.com/MlDenis/internal/gofermart/auth/cache"
	"github.com/MlDenis/internal/gofermart/storage"
)

// структура для наших хэндлеров, далее надо будет добавить возмонжо логер и тд
type HandlerOrderseDB struct {
	StorageOrders storage.InterfaceOrders
	DataJWT       *cache.DataJWT
}

func HandlerOrders(orders storage.InterfaceOrders, DataJWT *cache.DataJWT) *HandlerOrderseDB {
	return &HandlerOrderseDB{
		StorageOrders: orders,
		DataJWT:       DataJWT,
	}
}
