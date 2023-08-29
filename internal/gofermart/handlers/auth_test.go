package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MlDenis/internal/gofermart/auth"
	"github.com/MlDenis/internal/gofermart/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationUser_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/login",
		method:     http.MethodGet,
		user:       models.UserData{},
		statusCode: http.StatusMethodNotAllowed,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestAuthorizationUser_WrongContentType(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/login",
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

func TestAuthorizationUser_InvalidJSON(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       string
		user       models.UserData
		statusCode int
	}{
		request:    "/api/user/login",
		method:     http.MethodPost,
		body:       `{login: "test"`,
		user:       models.UserData{},
		statusCode: http.StatusInternalServerError,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestAuthorizationUser_Failed(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request: "/api/user/login",
		method:  http.MethodPost,
		user: models.UserData{
			Login:    "test",
			Password: "test",
		},
		statusCode: http.StatusUnauthorized,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.user)
	require.NoError(t, err)

	var pgErr pgconn.PgError

	ctx := context.TODO()

	data.user.PasswordHash = auth.HashPassword(data.user.Password)

	storeInterface.EXPECT().GetUser(ctx, &data.user).Return(&pgErr)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestAuthorizationUser_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		user       models.UserData
		statusCode int
	}{
		request: "/api/user/login",
		method:  http.MethodPost,
		user: models.UserData{
			Login:    "test",
			Password: "test",
		},
		statusCode: http.StatusOK,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.user)
	require.NoError(t, err)

	ctx := context.TODO()

	data.user.PasswordHash = auth.HashPassword(data.user.Password)

	storeInterface.EXPECT().GetUser(ctx, &data.user).Return(nil)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), data.user.Token, "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
	assert.NotNil(t, resp.Header.Get("models.HeaderHTTP"))
}
