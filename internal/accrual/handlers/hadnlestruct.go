package handlers

import (
	"github.com/MlDenis/internal/accrual/storage"
)

// структура для наших хэндлеров, далее надо будет добавить возмонжо логер и тд
type HandlerDB struct {
	Storage storage.DBInterfaceOrdersAccrual
}

func HandlerNew(s storage.DBInterfaceOrdersAccrual) *HandlerDB {
	return &HandlerDB{
		Storage: s,
	}
}
