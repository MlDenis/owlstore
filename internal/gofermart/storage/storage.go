package storage

import (
	"context"

	"github.com/MlDenis/internal/gofermart/models"
	"go.uber.org/zap"
)

type Interface interface {
	InterfaceUser
	InterfaceOrders
	InterfaceBalance
}
type InterfaceUser interface {
	RegisterUser(ctx context.Context, userData models.UserData) error
	GetUser(ctx context.Context, userData *models.UserData) error
	AuthorizationBalance(ctx context.Context, userlogin string) error
}

type InterfaceOrders interface {
	LoadOrderInDB(ctx context.Context, orderrData *models.Orders) error
	GetUserOrders(ctx context.Context, userlogin string) ([]models.OrdersOnly, error)
	GetAllOrders(ctx context.Context) ([]models.OrdersOnly, error)
	EditStatusAndAccrualOrder(ctx context.Context, status string, accrual, ordernumber int64) error
	EditBalanceAccrual(ctx context.Context, userlogin string, accrual int64) error
}
type InterfaceBalance interface {
	GetBalanceDB(ctx context.Context, userlogin string) (*models.ResponseBalance, error)
	LoadOrderInDB(ctx context.Context, orderrData *models.Orders) error
	EditBalanceWithdraw(ctx context.Context, userlogin string, sumwithdraw int64) error
	GetWithdrawalsDB(ctx context.Context, userlogin string) ([]models.WithdrawOrder, error)
}

func NewStorage(ctx context.Context, migratePath string, postgresDSN string, log *zap.Logger) (Interface, *PostgresDB, error) {

	DB, err := InitDB(postgresDSN, migratePath, log)
	if err != nil {
		log.Error("Error in initialization db", zap.Error(err))
		return nil, nil, err
	}
	return DB, DB, nil

}
