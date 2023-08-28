package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBalance_UnauthorizedUser(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/balance",
		method:     http.MethodGet,
		user:       models.UserData{},
		statusCode: http.StatusUnauthorized,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token)
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetBalance_EmptyUser(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       *models.UserData
		statusCode int
	}{
		request: "/api/user/balance",
		method:  http.MethodGet,
		user: &models.UserData{
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

	err := authUser(data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetBalanceDB(ctx, data.user.Login).Return(&models.ResponseBalance{
		AccrualSum:  0,
		WithdrawSum: 0,
	}, nil)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token)
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-type"))
}

func TestGetBalance_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       *models.UserData
		statusCode int
	}{
		request: "/api/user/balance",
		method:  http.MethodGet,
		user: &models.UserData{
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

	err := authUser(data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetBalanceDB(ctx, data.user.Login).Return(&models.ResponseBalance{
		AccrualSum:  12,
		WithdrawSum: 12,
	}, nil)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token)
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-type"))
}

func TestGetBalance_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       *models.UserData
		statusCode int
	}{
		request: "/api/user/balance",
		method:  http.MethodPost,
		user: &models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		statusCode: http.StatusMethodNotAllowed,
	}

	r, ctrl, handler, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	err := authUser(data.user, handler)
	require.NoError(t, err)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token)
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}
