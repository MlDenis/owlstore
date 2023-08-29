package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MlDenis/internal/gofermart/auth"
	"github.com/MlDenis/internal/gofermart/models"
	log "github.com/sirupsen/logrus"
)

// регистрируем нового пользователя
func (m *HandlerDB) RegisterNewUser(ctx context.Context) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			log.Error("got request with bad method", req.Method)
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			log.Error("wrong Content-Type", req.Header.Get("Content-Type"))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// десериализуем запрос в структуру модели
		log.Error("decoding request")
		//Создадим структуру пользователя, чтобы записать пользователя в бд
		jsonUsers := models.UserData{}
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&jsonUsers); err != nil {
			log.Error("cannot decode request JSON body", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonUsers.PasswordHash = auth.HashPassword(jsonUsers.Password)
		err := m.Storage.RegisterUser(ctx, jsonUsers)
		if err != nil {
			err := m.Storage.GetUser(ctx, &jsonUsers)
			if err == nil {
				log.Error("login busy")
				res.WriteHeader(http.StatusConflict)
			}
			log.Error("failed to register user", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//При регистрации создадим баланс для пользователя = 0
		err = m.Storage.AuthorizationBalance(ctx, jsonUsers.Login)
		if err != nil {
			log.Error("failed to register user", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//создаем токен для пользователя, который будет храниться в RAM
		jsonUsers.Token, err = auth.CreateJwtToken(jsonUsers.Login)
		if err != nil {
			log.Error("token not created")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		m.DataJWT.AddToken(&jsonUsers)
		res.Header().Add(models.HeaderHTTP, jsonUsers.Token)
		res.WriteHeader(http.StatusOK)

	}

}

// авторизация
func (m *HandlerDB) AuthorizationUser(ctx context.Context) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			log.Error("got request with bad method", req.Method)
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			log.Error("wrong Content-Type", req.Header.Get("Content-Type"))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// десериализуем запрос в структуру модели
		log.Error("decoding request")
		//Создадим структуру пользователя, чтобы записать токен и логин в него
		jsonUsers := &models.UserData{}
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(&jsonUsers); err != nil {
			log.Error("cannot decode request JSON body", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		//хэшируем пароль
		jsonUsers.PasswordHash = auth.HashPassword(jsonUsers.Password)
		err := m.Storage.GetUser(ctx, jsonUsers)
		if err != nil {
			log.Error("failed to autorization", err)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		//создаем токен для пользователя, который будет храниться в RAM
		jsonUsers.Token, err = auth.CreateJwtToken(jsonUsers.Login)
		if err != nil {
			log.Error("token not created")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		//При авторизации добавляем токен в RAM
		m.DataJWT.AddToken(jsonUsers)
		res.Header().Add(models.HeaderHTTP, jsonUsers.Token)
		res.WriteHeader(http.StatusOK)

	}

}
