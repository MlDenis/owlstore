package storage

import (
	"context"

	"github.com/MlDenis/internal/gofermart/models"
	log "github.com/sirupsen/logrus"
)

type DBInterface interface {
	RegisterUser(ctx context.Context, userData models.UserData) error
	GetUser(ctx context.Context, userData *models.UserData) error
	LoadOrderInDB(ctx context.Context, userData *models.Orders) error
	GetUserOrders(ctx context.Context, login string) ([]models.OrdersOnly, error)
	GetBalanceDB(ctx context.Context, userlogin string) (*models.ResponseBalance, error)
	AuthorizationBalance(ctx context.Context, userlogin string) error
	EditBalanceWithdraw(ctx context.Context, userlogin string, sumwithdraw int64) error
	GetWithdrawalsDB(ctx context.Context, login string) ([]models.WithdrawOrder, error)
	GetAllOrders(ctx context.Context) ([]models.OrdersOnly, error)
	EditStatusAndAccrualOrder(ctx context.Context, status string, accrual, ordernumber int64) error
	EditBalanceAccrual(ctx context.Context, userlogin string, accrual int64) error
}

func NewStorage(ctx context.Context, migratePath string, postgresDSN string) (DBInterface, *PostgresDB, error) {

	DB, err := InitDB(postgresDSN, migratePath)
	if err != nil {
		log.Error("Error in initialization db", (err))
		return nil, nil, err
	}
	return DB, DB, nil

}
