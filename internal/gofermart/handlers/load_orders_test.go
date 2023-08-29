package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/MlDenis/internal/gofermart/auth"
	"github.com/MlDenis/internal/gofermart/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOrderNumber_OrderAccepted(t *testing.T) {
	data := struct {
		name       string
		request    string
		method     string
		body       string
		user       models.UserData
		statusCode int
	}{
		name:    "new order has been accepted",
		request: "/api/user/orders",
		method:  http.MethodPost,
		body:    "9278923470",
		user: models.UserData{
			Login:        "test",
			Password:     "test",
			PasswordHash: auth.HashPassword("test"),
			Token:        "",
		},
		statusCode: http.StatusAccepted,
	}

	r, ctrl, handler, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	body := strings.NewReader(string(data.body))

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	order, err := makeOrder(data.user.Login, data.body)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().LoadOrderInDB(ctx, order)

	resp := testRequest(t, ts, data.method, data.request, body, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)

}

func makeOrder(login string, order string) (*models.Orders, error) {
	orderID, err := strconv.ParseInt(order, 10, 64)
	if err != nil {
		return &models.Orders{}, err
	}

	orderModel := &models.Orders{
		UserLogin:   login,
		OrderNumber: orderID,
	}

	return orderModel, nil
}

func authUser(user *models.UserData, handler *HandlerDB) error {
	token, err := auth.CreateJwtToken(user.Login)
	if err != nil {
		return err
	}

	user.Token = token

	handler.DataJWT.AddToken(user)

	return nil
}

func TestLoadOrderNumber_UnauthorizedUser(t *testing.T) {
	data := struct {
		name       string
		request    string
		method     string
		body       string
		user       models.UserData
		statusCode int
	}{
		name:       "unauthorized user",
		request:    "/api/user/orders",
		method:     http.MethodPost,
		body:       "",
		user:       models.UserData{},
		statusCode: http.StatusUnauthorized,
	}

	r, ctrl, _, _ := runTestServer(t)

	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	var token string

	resp := testRequest(t, ts, data.method, data.request, nil, token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestLoadOrderNumber_WrongOrderNumber(t *testing.T) {
	data := struct {
		name       string
		request    string
		method     string
		body       string
		user       models.UserData
		statusCode int
	}{
		name:    "wrong order number",
		request: "/api/user/orders",
		method:  http.MethodPost,
		body:    "test",
		user: models.UserData{
			Login:        "test",
			Password:     "test",
			PasswordHash: auth.HashPassword("test"),
			Token:        "",
		},
		statusCode: http.StatusBadRequest,
	}

	r, ctrl, handler, _ := runTestServer(t)

	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	body := strings.NewReader(string(data.body))

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	resp := testRequest(t, ts, data.method, data.request, body, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestLoadOrderNumber_InvalidOrderNumber(t *testing.T) {
	data := struct {
		name       string
		request    string
		method     string
		body       string
		user       models.UserData
		statusCode int
	}{
		name:    "invalid order number",
		request: "/api/user/orders",
		method:  http.MethodPost,
		body:    "123",
		user: models.UserData{
			Login:        "test",
			Password:     "test",
			PasswordHash: auth.HashPassword("test"),
			Token:        "",
		},
		statusCode: http.StatusUnprocessableEntity,
	}

	r, ctrl, handler, _ := runTestServer(t)

	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	body := strings.NewReader(string(data.body))

	err := authUser(&data.user, handler)
	require.NoError(t, err)

	resp := testRequest(t, ts, data.method, data.request, body, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestLoadOrderNumber_BadMethod(t *testing.T) {
	data := struct {
		name       string
		request    string
		method     string
		body       string
		user       *models.UserData
		statusCode int
	}{
		name:    "no orders",
		request: "/api/user/orders",
		method:  http.MethodPut,
		body:    "123",
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

	body := strings.NewReader(string(data.body))

	resp := testRequest(t, ts, data.method, data.request, body, data.user.Token, "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}
