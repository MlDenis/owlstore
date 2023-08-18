package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
)

func Router(ctx context.Context, newHandStruct *HandlerDB) chi.Router {
	r := chi.NewRouter()
	r.Post("/api/user/register", newHandStruct.RegisterNewUser(ctx))
	r.Post("/api/user/login", newHandStruct.AuthorizationUser(ctx))
	r.Post("/api/user/orders", newHandStruct.LoadOrderNumber(ctx))
	r.Get("/api/user/orders", newHandStruct.GetUserOrder(ctx))
	r.Get("/api/user/balance", newHandStruct.GetBalance(ctx))
	r.Post("/api/user/balance/withdraw", newHandStruct.WithdrawBalance(ctx))
	return r
}
