package storage

import (
	"context"

	"github.com/MlDenis/internal/accrual/models"
)

// Получение баланса пользователя
func (pgdb *PostgresDB) GetOrderFromOrdersAccrualDB(ctx context.Context, ordernumber int64) (*models.Order, error) {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return nil, err
	}

	ordersAccrual := &models.Order{}
	row := pgdb.pool.QueryRow(ctx, `SELECT ordernumber,statusorder,accrual FROM public.ordersaccrual WHERE ordernumber=$1`, ordernumber)
	err = row.Scan(&ordersAccrual.OrderNumber, &ordersAccrual.StatusOrder, &ordersAccrual.Accrual)
	if err != nil {

		tx.Rollback(ctx)
		return nil, err
	}

	return ordersAccrual, tx.Commit(ctx)
}

// Регистрация нового совершённого заказа
func (pgdb *PostgresDB) LoadOrderInOrdersAccrualDB(ctx context.Context, orderForRegister *models.OrderForRegister) error {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.ordersaccrual (ordernumber,statusorder,accrual,goods) VALUES ($1, $2, $3, $4)`,
		orderForRegister.OrderNumber, models.RegisteredOrder, 0, orderForRegister.Goods,
	)
	if err != nil {

		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// добавление aacrual
func (pgdb *PostgresDB) LoadAccrualStatusOrder(ctx context.Context, status string, ordernumber int64, accraul int64) error {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE public.ordersaccrual set accrual = $1, statusorder = $2 WHERE ordernumber=$3`,
		accraul, status, ordernumber,
	)
	if err != nil {

		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// получаем товары
func (pgdb *PostgresDB) GetAllOrdersAndGoods(ctx context.Context) ([]models.OrderForRegister, error) {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return nil, err
	}
	ordersGoods := []models.OrderForRegister{}
	rows, err := pgdb.pool.Query(ctx,
		`SELECT ordernumber,statusorder,goods FROM public.ordersaccrual`,
	)
	for rows.Next() {
		orderGoods := models.OrderForRegister{}
		err = rows.Scan(&orderGoods.OrderNumber, &orderGoods.StatusOrder, &orderGoods.Goods)
		if err != nil {

			tx.Rollback(ctx)
			return nil, err
		}
		ordersGoods = append(ordersGoods, orderGoods)
	}
	if err != nil {

		tx.Rollback(ctx)
		return nil, err
	}
	return ordersGoods, tx.Commit(ctx)
}
