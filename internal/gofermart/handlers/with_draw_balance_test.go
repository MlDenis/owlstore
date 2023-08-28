package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithdrawBalance_WrongContentType(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/balance/withdraw",
		method:     http.MethodPost,
		user:       models.UserData{},
		statusCode: http.StatusBadRequest,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/balance/withdraw",
		method:     http.MethodGet,
		user:       models.UserData{},
		statusCode: http.StatusMethodNotAllowed,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_UnauthorizedUser(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/balance/withdraw",
		method:     http.MethodPost,
		user:       models.UserData{},
		statusCode: http.StatusUnauthorized,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_InvalidOrderNumber(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		body       models.WithdrawOrder
		statusCode int
	}{
		request: "/api/user/balance/withdraw",
		method:  http.MethodPost,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		body: models.WithdrawOrder{
			Order: 1,
		},
		statusCode: http.StatusUnprocessableEntity,
	}

	r, ctrl, handler, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	err = authUser(&data.user, handler)
	require.NoError(t, err)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_InsufficientFunds(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		body       models.WithdrawOrder
		statusCode int
	}{
		request: "/api/user/balance/withdraw",
		method:  http.MethodPost,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		body: models.WithdrawOrder{
			Order: 2377225624,
			Sum:   12,
		},
		statusCode: http.StatusPaymentRequired,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	err = authUser(&data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetBalanceDB(ctx, data.user.Login).Return(&models.ResponseBalance{
		AccrualSum:  0,
		WithdrawSum: 0,
	}, nil)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_UniqueViolation(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		body       models.WithdrawOrder
		statusCode int
	}{
		request: "/api/user/balance/withdraw",
		method:  http.MethodPost,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		body: models.WithdrawOrder{
			Order: 2377225624,
			Sum:   12,
		},
		statusCode: http.StatusConflict,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	err = authUser(&data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetBalanceDB(ctx, data.user.Login).Return(&models.ResponseBalance{
		AccrualSum:  100,
		WithdrawSum: 100,
	}, nil)

	order := &models.Orders{}
	order.OrderNumber = data.body.Order
	order.UserLogin = data.user.Login
	order.StatusOrder = models.WithdrawEnd
	order.Withdraw = data.body.Sum

	var pgErr pgconn.PgError
	pgErr.Code = pkg.UniqueViolationCode

	storeInterface.EXPECT().LoadOrderInDB(ctx, order).Return(&pgErr) // возвращается conflict

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_EditBalanceWithdraw(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		body       models.WithdrawOrder
		statusCode int
	}{
		request: "/api/user/balance/withdraw",
		method:  http.MethodPost,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		body: models.WithdrawOrder{
			Order: 2377225624,
			Sum:   12,
		},
		statusCode: http.StatusOK,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	err = authUser(&data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetBalanceDB(ctx, data.user.Login).Return(&models.ResponseBalance{
		AccrualSum:  100,
		WithdrawSum: 100,
	}, nil)

	order := &models.Orders{}
	order.OrderNumber = data.body.Order
	order.UserLogin = data.user.Login
	order.StatusOrder = models.WithdrawEnd
	order.Withdraw = data.body.Sum

	storeInterface.EXPECT().LoadOrderInDB(ctx, order).Return(nil)

	storeInterface.EXPECT().EditBalanceWithdraw(ctx, data.user.Login, order.Withdraw).Return(nil)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestWithdrawBalance_InvalidJSON(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		body       string
		statusCode int
	}{
		request:    "/api/user/balance/withdraw",
		method:     http.MethodPost,
		user:       models.UserData{},
		body:       `{"abc": 1`,
		statusCode: http.StatusInternalServerError,
	}

	r, ctrl, handler, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	body := strings.NewReader(data.body)

	resp := testRequest(t, ts, data.method, data.request, body, data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}
