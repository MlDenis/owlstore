package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MlDenis/internal/accrual/models"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterNewOrder_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders",
		method:     http.MethodGet,
		statusCode: http.StatusMethodNotAllowed,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterNewOrder_BadContentType(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, "", "test/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterNewOrder_InvalidJSON(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       string
		statusCode int
	}{
		request:    "/api/orders",
		method:     http.MethodPost,
		body:       `{order_number: 1`,
		statusCode: http.StatusInternalServerError,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterNewOrder_InvalidORderNumber(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       models.OrderForRegister
		statusCode int
	}{
		request: "/api/orders",
		method:  http.MethodPost,
		body: models.OrderForRegister{
			OrderNumber: 1,
		},
		statusCode: http.StatusUnprocessableEntity,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterNewOrder_LoadOrderConflict(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       models.OrderForRegister
		statusCode int
	}{
		request: "/api/orders",
		method:  http.MethodPost,
		body: models.OrderForRegister{
			OrderNumber: 9278923470,
			Goods: []models.Goods{
				{
					Description: "Чайник Bork",
					Price:       100,
				},
			},
		},
		statusCode: http.StatusConflict,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	pgErr := &pgconn.PgError{}
	pgErr.Code = pkg.UniqueViolationCode

	ctx := context.TODO()

	storeInterface.EXPECT().LoadOrderInOrdersAccrualDB(ctx, &data.body).Return(pgErr)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterNewOrder_LoadOrderDBError(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       models.OrderForRegister
		statusCode int
	}{
		request: "/api/orders",
		method:  http.MethodPost,
		body: models.OrderForRegister{
			OrderNumber: 9278923470,
			Goods: []models.Goods{
				{
					Description: "Чайник Bork",
					Price:       100,
				},
			},
		},
		statusCode: http.StatusBadRequest,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	pgErr := &pgconn.PgError{}

	ctx := context.TODO()

	storeInterface.EXPECT().LoadOrderInOrdersAccrualDB(ctx, &data.body).Return(pgErr)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterNewOrder_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       models.OrderForRegister
		statusCode int
	}{
		request: "/api/orders",
		method:  http.MethodPost,
		body: models.OrderForRegister{
			OrderNumber: 9278923470,
			Goods: []models.Goods{
				{
					Description: "Чайник Bork",
					Price:       100,
				},
			},
		},
		statusCode: http.StatusAccepted,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	bodyJSON, err := json.Marshal(data.body)
	require.NoError(t, err)

	ctx := context.TODO()

	storeInterface.EXPECT().LoadOrderInOrdersAccrualDB(ctx, &data.body).Return(nil)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(bodyJSON), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}
