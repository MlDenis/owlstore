package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MlDenis/internal/accrual/models"
	"github.com/MlDenis/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestGetOrder_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders/9278923470",
		method:     http.MethodPost,
		statusCode: http.StatusMethodNotAllowed,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, "", "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetOrder_WrongOrderNumber(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders/ssss",
		method:     http.MethodGet,
		statusCode: http.StatusBadRequest,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, "", "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetOrder_InvalidOrderNumber(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders/123",
		method:     http.MethodGet,
		statusCode: http.StatusUnprocessableEntity,
	}

	r, ctrl, _, _ := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp := testRequest(t, ts, data.method, data.request, nil, "", "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetOrder_NoContent(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders/9278923470",
		method:     http.MethodGet,
		statusCode: http.StatusNoContent,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.TODO()

	storeInterface.EXPECT().GetOrderFromOrdersAccrualDB(ctx, int64(9278923470)).Return(&models.Order{}, pkg.NoOrders)

	resp := testRequest(t, ts, data.method, data.request, nil, "", "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetOrder_DBError(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders/9278923470",
		method:     http.MethodGet,
		statusCode: http.StatusInternalServerError,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.TODO()

	pgErr := pgconn.PgError{}

	storeInterface.EXPECT().GetOrderFromOrdersAccrualDB(ctx, int64(9278923470)).Return(&models.Order{}, &pgErr)

	resp := testRequest(t, ts, data.method, data.request, nil, "", "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestGetOrder_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/orders/9278923470",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.TODO()

	order := &models.Order{
		OrderNumber: int64(9278923470),
		StatusOrder: "PROCESSING",
		Accrual:     15,
	}

	storeInterface.EXPECT().GetOrderFromOrdersAccrualDB(ctx, int64(9278923470)).Return(order, nil)

	resp := testRequest(t, ts, data.method, data.request, nil, "", "text/plain")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-type"))
}
