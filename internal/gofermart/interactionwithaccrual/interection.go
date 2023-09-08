package interactionwithaccrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/internal/gofermart/storage"
	"go.uber.org/zap"
)

// WorkerPool принимает канал данных, порождает 10 горутин
func WorkerPool(ctx context.Context, s storage.Interface, rateLimit int, url string, log *zap.Logger) {
	jobs := make(chan models.OrdersOnly, rateLimit)
	// g := new(errgroup.Group)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(10) * time.Second):
			for w := 1; w <= rateLimit; w++ {
				go GetAccrualAndStatus(ctx, jobs, s, url, log)
			}
			go OrdersGoodsGorutine(ctx, jobs, s, log)
		}

	}

}

func OrdersGoodsGorutine(ctx context.Context, ordersChan chan models.OrdersOnly, s storage.Interface, log *zap.Logger) {
	rewards, err := s.GetAllOrders(ctx)
	if err != nil {
		log.Error("error in get rewards from db: ", zap.Error(err))
		ordersChan <- models.OrdersOnly{OrderNumber: 0}
		return

	}

	for _, order := range rewards {
		err = s.EditStatusAndAccrualOrder(ctx, models.ProcessingOrder, 0, order.OrderNumber)
		if err != nil {
			log.Error("error in add accrual in db: ", zap.Error(err))
			return
		}
		ordersChan <- order
	}

}

func GetAccrualAndStatus(ctx context.Context, ordersChan chan models.OrdersOnly, s storage.Interface, url string, log *zap.Logger) {

	order := <-ordersChan
	if order.OrderNumber == 0 {
		return
	}
	orderNumberToString := strconv.FormatInt(order.OrderNumber, 10)
	urlGet := "http://" + url + "/api/orders/" + string(orderNumberToString)
	resp, err := http.Get(urlGet)
	if err != nil {
		log.Error("connection refuser: ", zap.Error(err))
		return
	}

	if resp.Status != strconv.Itoa(http.StatusAccepted) {
		errorStatus := fmt.Sprintf("order with number %s does not exist", orderNumberToString)
		log.Error(errorStatus)
		return
	}
	orderResp := &models.OrderResp{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(orderResp); err != nil {
		log.Error("cannot decode request JSON body: ", zap.Error(err))
		return
	}
	err = s.EditStatusAndAccrualOrder(ctx, models.ProcessedOrder, orderResp.Accrual, order.OrderNumber)
	if err != nil {
		log.Error("error in add accrual in db: ", zap.Error(err))
		return
	}
	err = s.EditBalanceAccrual(ctx, order.UserLogin, orderResp.Accrual)
	if err != nil {
		log.Error("error in add accrual in db: ", zap.Error(err))
		return
	}
}
