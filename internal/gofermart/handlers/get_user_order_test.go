package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MlDenis/internal/gofermart/models"
	"github.com/MlDenis/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserOrder_UnauthorizedUser(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/orders",
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

func TestGetUserOrder_NoOrders(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       *models.UserData
		statusCode int
	}{
		request: "/api/user/orders",
		method:  http.MethodGet,
		user: &models.UserData{
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

	err := authUser(data.user, handler)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().GetUserOrders(ctx, data.user).Return([]models.Orders{}, pkg.NoOrders)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token)
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetUserOrder_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       *models.UserData
		statusCode int
	}{
		request: "/api/user/orders",
		method:  http.MethodPost,
		user: &models.UserData{
			Login:        "test",
			Password:     "",
			PasswordHash: "",
			Token:        "",
		},
		statusCode: http.StatusBadRequest,
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

func TestGetUserOrder_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       *models.UserData
		statusCode int
	}{
		request: "/api/user/orders",
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

	storeInterface.EXPECT().GetUserOrders(ctx, data.user).Return([]models.Orders{
		{
			UserLogin:   data.user.Login,
			OrderNumber: 9278923470,
			OrderDate:   time.Now(),
			StatusOrder: "PROCESSING",
		},
	}, nil)

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token)
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-type"))
}
