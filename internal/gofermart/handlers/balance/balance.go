package balance

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/internal/luna"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// баланс пользователя
func (m *HandlerBalanceDB) GetBalance(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			log.Error("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		jsonUsers := &models.UserData{}
		jsonUsers.Token = req.Header.Get(models.HeaderHTTP)
		err := m.DataJWT.GetToken(jsonUsers)
		if err != nil {
			log.Error("user not authenticated: ", zap.Error(err))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ResponseBalance, err := m.StorageBalance.GetBalanceDB(ctx, jsonUsers.Login)
		if err != nil {
			log.Error("error in add orders in db: ", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		response, err := json.Marshal(ResponseBalance)
		if err != nil {
			log.Error("cannot marshal to json: ", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(response)
		return
	}
}

// Запрос на списание средств
func (m *HandlerBalanceDB) WithdrawBalance(ctx context.Context, log *zap.Logger) http.HandlerFunc {
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
		//Создадим структуру пользователя, чтобы записать токен и логин в него
		jsonUsers := &models.UserData{}
		jsonUsers.Token = req.Header.Get(models.HeaderHTTP)
		err := m.DataJWT.GetToken(jsonUsers)
		if err != nil {
			log.Error("user not authenticated: ", zap.Error(err))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		wisthdrawSum := &models.WithdrawOrder{}
		// десериализуем запрос в структуру модели
		log.Error("decoding request")
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&wisthdrawSum); err != nil {
			log.Error("cannot decode request JSON body", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		//проверям заказ через алгоритм луна
		validNumber := luna.Valid(wisthdrawSum.Order)
		if !validNumber {
			log.Error("invalid order number")
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		//Проверям баланс
		ResponseBalance, err := m.StorageBalance.GetBalanceDB(ctx, jsonUsers.Login)
		if err != nil {
			log.Error("error in add orders in db: ", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		//проверям хватает ли средств для списания
		if ResponseBalance.AccrualSum < wisthdrawSum.Sum {
			log.Error("insufficient funds to write off: ", zap.Error(err))
			res.WriteHeader(http.StatusPaymentRequired)
			return
		}
		//Создадим структуру заказов, чтобы записать их в бд
		jsonOrders := &models.Orders{}
		jsonOrders.OrderNumber = wisthdrawSum.Order
		jsonOrders.UserLogin = jsonUsers.Login
		jsonOrders.StatusOrder = models.WithdrawEnd
		jsonOrders.Withdraw = wisthdrawSum.Sum
		err = m.StorageBalance.LoadOrderInDB(ctx, jsonOrders)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				//Проверяем,не загружали ли этот заказ ранее(точно ли это ноывй заказ)
				if pgErr.Code == pkg.UniqueViolationCode {
					log.Error("cannot be debited for this order because this order has already been uploaded")
					res.WriteHeader(http.StatusConflict)
					return
				}
				log.Error("error in add orders in db: ", zap.Error(err))
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Error("error in add orders in db: ", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//Изменим баланс после списания
		err = m.StorageBalance.EditBalanceWithdraw(ctx, jsonOrders.UserLogin, wisthdrawSum.Sum)
		if err != nil {
			log.Error("balance change error: ", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		return
	}
}

// Получение информации о выводе средств
func (m *HandlerBalanceDB) GetWithdrawals(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			log.Error("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		userData := &models.UserData{}
		userData.Token = req.Header.Get(models.HeaderHTTP)

		err := m.DataJWT.GetToken(userData)
		if err != nil {
			log.Error("user not authenticated: ", zap.Error(err))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		withdrawals, err := m.StorageBalance.GetWithdrawalsDB(ctx, userData.Login)
		if err != nil {
			if errors.Is(err, pkg.NoOrders) {
				res.WriteHeader(http.StatusNoContent)
				return
			}
			log.Error("cannot get user's withdrawals balance: ", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		withdrawalsJson, err := json.Marshal(withdrawals)
		if err != nil {
			log.Error("cannot make json withdrawals balance: ", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-type", "application/json")

		_, err = res.Write(withdrawalsJson)
		if err != nil {
			log.Error("cannot orders json: ", zap.Error(err))
			return
		}
		res.WriteHeader(http.StatusOK)
		return
	}
}
