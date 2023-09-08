package handlers

import (
	"context"

	"github.com/MlDenis/internal/gofermart/handlers/balance"
	"github.com/MlDenis/internal/gofermart/handlers/order"
	"github.com/MlDenis/internal/gofermart/handlers/users"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Router(ctx context.Context, log *zap.Logger, newHandStruct *HandlerDB) chi.Router {
	Balance := balance.HandlerBalance(newHandStruct.Storage, newHandStruct.DataJWT)
	Users := users.HandlerUsers(newHandStruct.Storage, newHandStruct.DataJWT)
	Orders := order.HandlerOrders(newHandStruct.Storage, newHandStruct.DataJWT)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/user/register", Users.RegisterNewUser(ctx, log))
	r.Post("/api/user/login", Users.AuthorizationUser(ctx, log))
	r.Post("/api/user/orders", Orders.LoadOrderNumber(ctx, log))
	r.Get("/api/user/orders", Orders.GetUserOrder(ctx, log))
	r.Get("/api/user/balance", Balance.GetBalance(ctx, log))
	r.Post("/api/user/balance/withdraw", Balance.WithdrawBalance(ctx, log))
	r.Get("/api/user/withdrawals", Balance.GetWithdrawals(ctx, log))
	return r
}
