package order

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/internal/luna"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// авторизация
func (m *HandlerOrderseDB) LoadOrderNumber(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			log.Error("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if req.Header.Get("Content-Type") != "text/plain" {
			log.Error("wrong Content-Type: ", zap.String("method", req.Header.Get("Content-Type")))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//Создадим структуру пользователя, чтобы записать токен и логин в него
		jsonUsers := &models.UserData{}
		//Создадим структуру заказов, чтобы записать их в бд
		jsonOrders := &models.Orders{}
		jsonUsers.Token = req.Header.Get(models.HeaderHTTP)
		err := m.DataJWT.GetToken(jsonUsers)
		if err != nil {
			log.Error("user not authenticated: ", zap.Error(err))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		number, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("error in read request body: ", zap.Error(err))
			res.WriteHeader(http.StatusBadRequest)
		}
		//проверям заказ через алгоритм луна
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
		jsonOrders.OrderNumber = orderID
		jsonOrders.UserLogin = jsonUsers.Login
		err = m.StorageOrders.LoadOrderInDB(ctx, jsonOrders)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pkg.UniqueViolationCode {
					log.Error("the order number has already been uploaded by the user")
					res.WriteHeader(http.StatusOK)
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
		log.Error("new order number accepted for processing")
		res.WriteHeader(http.StatusAccepted)
		return
	}

}

func (m *HandlerOrderseDB) GetUserOrder(ctx context.Context, log *zap.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			log.Error("got request with bad method", zap.String("method", req.Method))
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		//Проверяем токен
		userData := &models.UserData{}
		userData.Token = req.Header.Get(models.HeaderHTTP)

		err := m.DataJWT.GetToken(userData)
		if err != nil {
			log.Error("user not authenticated: ", zap.Error(err))
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		orders, err := m.StorageOrders.GetUserOrders(ctx, userData.Login)
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

		_, err = res.Write(ordersJson)
		if err != nil {
			log.Error("cannot orders json: ", zap.Error(err))
			return
		}
	}
}
