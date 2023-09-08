package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/MlDenis/internal/accrual/models"
	"github.com/MlDenis/internal/luna"
	"github.com/MlDenis/pkg"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// Регаем новый заказ
func (m *HandlerDB) RegisterNewOrder(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			log.Error("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			log.Error("wrong Content-Type", zap.String("method", req.Header.Get("Content-Type")))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonOrder := &models.OrderForRegister{}
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(jsonOrder); err != nil {
			log.Error("cannot decode request JSON body", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		validNumber := luna.Valid(jsonOrder.OrderNumber)
		if !validNumber {
			log.Error("invalid order number")
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		//Записываем новый заказ в бд
		err := m.Storage.LoadOrderInOrdersAccrualDB(ctx, jsonOrder)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pkg.UniqueViolationCode {
					log.Error("the order number has already been uploaded by the user")
					res.WriteHeader(http.StatusConflict)
					return
				}
			}
			log.Error("error in add orders in db", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusAccepted)
		return
	}

}

// Получаем accraul заказа и его статус
func (m *HandlerDB) GetOrder(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			log.Error("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		number := chi.URLParam(req, "number")
		//Проверяем на алгоритм луна
		orderID, err := strconv.ParseInt(string(number), 10, 64)
		if err != nil {
			log.Error("wrong order number:", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		validNumber := luna.Valid(orderID)
		if !validNumber {
			log.Error("invalid order number")
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		//Смотрим данные в бд
		orders, err := m.Storage.GetOrderFromOrdersAccrualDB(ctx, orderID)
		if err != nil {
			if errors.Is(err, pkg.NoOrders) {
				res.WriteHeader(http.StatusNoContent)
				return
			}
			log.Error("cannot get user's orders: ", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		ordersJson, err := json.Marshal(orders)
		if err != nil {
			log.Error("cannot make json orders: ", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(ordersJson)
		return
	}
}
