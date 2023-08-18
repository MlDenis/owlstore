package storage

import (
	"context"

	"github.com/MlDenis/internal/gofermart/models"
	log "github.com/sirupsen/logrus"
)

// Получение баланса пользователя
func (pgdb *PostgresDB) GetBalanceDB(ctx context.Context, userlogin string) (*models.ResponseBalance, error) {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var (
		AccrualSum  int64
		WithdrawSum int64
	)
	ResponseBalance := &models.ResponseBalance{}
	row := pgdb.pool.QueryRow(ctx, `SELECT sumaccrual,sumwithdraw FROM public.balance WHERE userlogin=$1`, userlogin)
	err = row.Scan(&AccrualSum, &WithdrawSum)
	if err != nil {
		log.Error(err)
		tx.Rollback(ctx)
		return nil, err
	}
	ResponseBalance.AccrualSum = AccrualSum
	ResponseBalance.WithdrawSum = WithdrawSum

	return ResponseBalance, tx.Commit(ctx)
}

// При авторизации сразу заполняем баланс пользователя по нулям
func (pgdb *PostgresDB) AuthorizationBalance(ctx context.Context, userlogin string) error {

	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO public.balance (userlogin,sumaccrual,sumwithdraw) VALUES ($1, $2,$3)`, userlogin, models.BalanceAuthAccrualWithdraw, models.BalanceAuthAccrualWithdraw)
	if err != nil {
		log.Error(err)
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// Меняем баланс при списании
func (pgdb *PostgresDB) EditBalanceWithdraw(ctx context.Context, userlogin string, sumwithdraw int64) error {

	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO public.balance (userlogin, sumaccrual, sumwithdraw)
		VALUES ($1, $2, $3)
		ON CONFLICT (userlogin) DO UPDATE
		SET sumaccrual = public.balance.sumaccrual - EXCLUDED.sumaccrual,
		sumwithdraw = public.balance.sumwithdraw + EXCLUDED.sumwithdraw `,
		userlogin, sumwithdraw, sumwithdraw)
	if err != nil {
		log.Error(err)
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
