package storage

import (
	"context"

	"github.com/MlDenis/internal/accrual/models"
)

func (pgdb *PostgresDB) RegisterInfoInDB(ctx context.Context, goods *models.Reward) error {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO public.rewards (match,reward,reward_type) VALUES ($1, $2, $3)`,
		goods.Match, goods.Reward, goods.RewardType,
	)
	if err != nil {

		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (pgdb *PostgresDB) GetAllRewards(ctx context.Context) ([]models.Reward, error) {
	tx, err := pgdb.pool.Begin(ctx)
	if err != nil {

		return nil, err
	}
	rewardArr := []models.Reward{}

	rows, err := pgdb.pool.Query(ctx, `SELECT match,reward,reward_type FROM public.rewards`)

	for rows.Next() {
		reward := models.Reward{}
		err = rows.Scan(&reward.Match, &reward.Reward, &reward.RewardType)
		if err != nil {

			tx.Rollback(ctx)
			return nil, err
		}
		rewardArr = append(rewardArr, reward)
	}
	if err != nil {

		tx.Rollback(ctx)
		return nil, err
	}

	return rewardArr, tx.Commit(ctx)
}
