package storage

import (
	"context"

	"github.com/MlDenis/internal/accrual/models"
	log "github.com/sirupsen/logrus"
)

type DBInterfaceOrdersAccrual interface {
	GetOrderFromOrdersAccrualDB(ctx context.Context, ordernumber int64) (*models.Order, error)
	LoadOrderInOrdersAccrualDB(ctx context.Context, order *models.OrderForRegister) error
	RegisterInfoInDB(ctx context.Context, goods *models.Reward) error
	// AddGoods(ctx context.Context, orderForRegister *models.OrderForRegister) error
	// GetAllGoods(ctx context.Context, orders *models.OrderForRegister) ([]models.GoodsWithReward, error)
	LoadAccrualStatusOrder(ctx context.Context, status string, ordernumber, accraul int64) error
	GetAllOrdersAndGoods(ctx context.Context) ([]models.OrderForRegister, error)
	GetAllRewards(ctx context.Context) ([]models.Reward, error)
}

func NewStorage(ctx context.Context, migratePath string, postgresDSN string) (DBInterfaceOrdersAccrual, *PostgresDB, error) {

	DB, err := InitDB(postgresDSN, migratePath)
	if err != nil {
		log.Error("Error in initialization db", (err))
		return nil, nil, err
	}
	return DB, DB, nil

}
