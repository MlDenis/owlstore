package handlers

import (
	"context"
	"net/http"

	"github.com/MlDenis/internal/accrual/models"
	"github.com/MlDenis/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"
)

// сделаем апи после того как проверим что все работает
func Router(ctx context.Context, log *zap.Logger, newHandStruct *HandlerDB) chi.Router {
	r := chi.NewRouter()
	//лимит для кол-ва запросов
	r.Use(httprate.Limit(
		models.RateLimit,
		models.TimeLimit,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "some specific response here", http.StatusTooManyRequests)
		}),
	))
	r.Use((logger.WithLogging(log)))
	r.Post("/api/orders", newHandStruct.RegisterNewOrder(ctx, log))
	r.Post("/api/goods", newHandStruct.RegisterInfoReward(ctx, log))
	r.Get("/api/orders/{number}", newHandStruct.GetOrder(ctx, log))
	return r
}
