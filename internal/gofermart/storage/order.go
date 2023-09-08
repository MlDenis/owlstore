package storage

import (
	"context"
	"time"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/pkg"
)

// записываем заказы пользователя
func (pgdb *PostgresDB) LoadOrderInDB(ctx context.Context, orders *models.Orders) error {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return err
	}
	orders.OrderDate = time.Now()
	//Проверяем, если это заказ на списание, то добавляем статус models.WithdrawEnd, а если заказ новый, который ждет начисление, то models.NewOrder
	if orders.StatusOrder == models.WithdrawEnd {
		_, err = tx.Exec(ctx,
			`INSERT INTO public.orders (ordernumber,userlogin,orderdate,statusorder,withdraw) VALUES ($1, $2,$3, $4, $5)`,
			orders.OrderNumber, orders.UserLogin, orders.OrderDate, orders.StatusOrder, orders.Withdraw,
		)
		if err != nil {

			tx.Rollback(ctx)
			return err
		}
		return tx.Commit(ctx)
	}
	orders.StatusOrder = models.NewOrder
	_, err = tx.Exec(ctx,
		`INSERT INTO public.orders (ordernumber,userlogin,orderdate,statusorder) VALUES ($1, $2,$3, $4)`,
		orders.OrderNumber, orders.UserLogin, orders.OrderDate, orders.StatusOrder,
	)
	if err != nil {

		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)

}

func (pgdb *PostgresDB) GetUserOrders(ctx context.Context, userlogin string) ([]models.OrdersOnly, error) {
	orders := []models.OrdersOnly{}
	rows, err := pgdb.pool.Query(ctx, `SELECT ordernumber, orderdate, statusorder FROM public.orders WHERE userlogin = $1`, userlogin) // дописать accrual, withdraw когда сделаем систему
	if err != nil {
		return orders, err
	}

	for rows.Next() {
		order := models.OrdersOnly{}
		err := rows.Scan(&order.OrderNumber, &order.OrderDate, &order.StatusOrder)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	defer rows.Close()

	if len(orders) == 0 {
		return orders, pkg.NoOrders
	}

	return orders, nil
}

// Получение информации о выводе средств из бд
func (pgdb *PostgresDB) GetWithdrawalsDB(ctx context.Context, userlogin string) ([]models.WithdrawOrder, error) {
	withdrawals := []models.WithdrawOrder{}
	rows, err := pgdb.pool.Query(ctx, `SELECT ordernumber, withdraw, orderdate FROM public.orders WHERE userlogin = $1 and statusorder=$2`, userlogin, models.WithdrawEnd) // дописать accrual, withdraw когда сделаем систему
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		withdraw := models.WithdrawOrder{}
		err := rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	defer rows.Close()

	if len(withdrawals) == 0 {
		return nil, pkg.NoOrders
	}

	return withdrawals, nil
}

func (pgdb *PostgresDB) GetAllOrders(ctx context.Context) ([]models.OrdersOnly, error) {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return nil, err
	}
	orders := []models.OrdersOnly{}
	rows, err := pgdb.pool.Query(ctx, `SELECT ordernumber, orderdate, statusorder,userlogin FROM public.orders WHERE statusorder=$1 or statusorder=$2`, models.NewOrder, models.ProcessingOrder) // дописать accrual, withdraw когда сделаем систему
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		order := models.OrdersOnly{}
		err := rows.Scan(&order.OrderNumber, &order.OrderDate, &order.StatusOrder, &order.UserLogin)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return orders, pkg.NoOrders
	}
	if err != nil {

		tx.Rollback(ctx)
		return nil, err
	}

	return orders, tx.Commit(ctx)

}

func (pgdb *PostgresDB) EditStatusAndAccrualOrder(ctx context.Context, status string, accrual, ordernumber int64) error {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE public.orders set accrual = $1, statusorder = $2 WHERE ordernumber=$3`,
		accrual, status, ordernumber,
	)
	if err != nil {

		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}
