package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MlDenis/internal/accrual/models"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// Регистрация информации о вознаграждении за товар
func (m *HandlerDB) RegisterInfoReward(ctx context.Context, log *zap.Logger) http.HandlerFunc {
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

		jsonGoods := &models.Reward{}
		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(jsonGoods); err != nil {
			log.Error("cannot decode request JSON body", zap.Error(err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		err := m.Storage.RegisterInfoInDB(ctx, jsonGoods)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pkg.UniqueViolationCode {
					log.Error("the order number has already been uploaded by the user")
					res.WriteHeader(http.StatusConflict)
					return
				}
			}
			log.Error("error in add in db: ", zap.Error(err))
			res.WriteHeader(http.StatusConflict)
			return
		}

		res.WriteHeader(http.StatusAccepted)
		return
	}
}
