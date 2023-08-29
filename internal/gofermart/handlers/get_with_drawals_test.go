package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWithDrawals_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/withdrawals",
		method:     http.MethodPost,
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

func TestGetWithDrawals_UnauthorizedUser(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/withdrawals",
		method:     http.MethodGet,
		user:       models.UserData{},
		statusCode: http.StatusUnauthorized,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetWithDrawals_NoContent(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request: "/api/user/withdrawals",
		method:  http.MethodGet,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		statusCode: http.StatusNoContent,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetWithdrawalsDB(ctx, data.user.Login).Return([]models.WithdrawOrder{}, pkg.NoOrders)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetWithDrawals_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request: "/api/user/withdrawals",
		method:  http.MethodGet,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		statusCode: http.StatusOK,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetWithdrawalsDB(ctx, data.user.Login).Return([]models.WithdrawOrder{}, nil)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetWithDrawals_DBError(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request: "/api/user/withdrawals",
		method:  http.MethodGet,
		user: models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		statusCode: http.StatusInternalServerError,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	var dbErr *pgconn.PgError

	storeInterface.EXPECT().GetWithdrawalsDB(ctx, data.user.Login).Return([]models.WithdrawOrder{}, dbErr)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}
