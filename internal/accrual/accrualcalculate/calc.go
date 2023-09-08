package accrualcalculate

import (
	"context"
	"regexp"
	"time"

	"github.com/MlDenis/internal/accrual/models"
	"github.com/MlDenis/internal/accrual/storage"
	"go.uber.org/zap"
)

// WorkerPool принимает канал данных, порождает 10 горутин
func WorkerPool(ctx context.Context, s storage.DBInterfaceOrdersAccrual, rateLimit int, log *zap.Logger) {
	jobs := make(chan models.GoodsWithReward, rateLimit)
	// g := new(errgroup.Group)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(100) * time.Second):
			for w := 1; w <= rateLimit; w++ {
				go resultAccrual(ctx, jobs, s, log)
			}

			go OrdersGoodsGorutine(ctx, jobs, s, log)
		}

	}

}

func OrdersGoodsGorutine(ctx context.Context, ordersChan chan models.GoodsWithReward, s storage.DBInterfaceOrdersAccrual, log *zap.Logger) {

	// select {
	// case <-ctx.Done():
	// 	return
	// case <-time.After(time.Duration(5) * time.Second):
	goodsWithReward := models.GoodsWithReward{}
	rewards, err := s.GetAllRewards(ctx)
	if err != nil {
		log.Error("error in get rewards from db: ", zap.Error(err))
		goodsWithReward.OrderNumber = 0
		return
	}
	ordersWithGoods, err := s.GetAllOrdersAndGoods(ctx)
	if err != nil {
		log.Error("error in get orders from db: ", zap.Error(err))
		goodsWithReward.OrderNumber = 0
		return
	}
	goodsWithReward.Reward = rewards
	for i := 0; i < len(ordersWithGoods); i++ {
		if ordersWithGoods[i].StatusOrder != models.ProcessedOrder || ordersWithGoods[i].StatusOrder != models.InvalidOrder {
			err := s.LoadAccrualStatusOrder(ctx, models.ProcessingOrder, goodsWithReward.OrderForRegister.OrderNumber, 0)
			if err != nil {
				log.Error("error in add orders from db: ", zap.Error(err))
				return
			}
			goodsWithReward.OrderForRegister = ordersWithGoods[i]
			ordersChan <- goodsWithReward
		}
	}
}

func resultAccrual(ctx context.Context, orderAndRewardChan chan models.GoodsWithReward, s storage.DBInterfaceOrdersAccrual, log *zap.Logger) {

	var accraulSum int64 = 0
	orderAndReward := <-orderAndRewardChan
	if orderAndReward.OrderNumber == 0 {
		return
	}
	for _, reward := range orderAndReward.Reward {
		for _, goods := range orderAndReward.Goods {
			matched, err := regexp.MatchString(reward.Match, goods.Description)
			if err != nil {
				err := s.LoadAccrualStatusOrder(ctx, models.InvalidOrder, orderAndReward.OrderNumber, 0)
				if err != nil {
					log.Error("error in add orders from db: ", zap.Error(err))
					return
				}
				log.Error("error in get orders from db: ", zap.Error(err))
				return
			}
			if matched {
				accrualOne := accrualCalculate(reward.RewardType, goods.Price, reward.Reward)
				accraulSum += accrualOne
			}
		}
	}
	err := s.LoadAccrualStatusOrder(ctx, models.ProcessedOrder, orderAndReward.OrderNumber, accraulSum)
	if err != nil {
		log.Error("error in add orders from db: ", zap.Error(err))
		return
	}

}

func accrualCalculate(rewardType string, price, Reward int64) int64 {
	accrual := 0.00
	if rewardType == models.RewardTypeDefault {
		rewardToFloat64 := float64(Reward)
		priceFloat64 := float64(price)
		sum := rewardToFloat64 * priceFloat64 / 100
		accrual += sum
	} else {
		return Reward
	}
	accrualToInt64 := int64(accrual)
	return accrualToInt64
}
