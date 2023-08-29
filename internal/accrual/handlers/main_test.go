package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MlDenis/internal/accrual/accrualcalculate"
	mock_storage "github.com/MlDenis/internal/accrual/storage/mocks"
	"github.com/MlDenis/internal/gofermart/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader, token string, contentType string) *http.Response {

	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Close = true
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "PostmanRuntime/7.32.3")
	req.Header.Add("Content-Type", contentType)
	req.Header.Add(models.HeaderHTTP, token)

	require.NoError(t, err)

	ts.Client()

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func runTestServer(t *testing.T) (chi.Router, *gomock.Controller, *HandlerDB, *mock_storage.MockDBInterfaceOrdersAccrual) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	memStorageInterface := mock_storage.NewMockDBInterfaceOrdersAccrual(ctrl)

	newHandStruct := HandlerNew(memStorageInterface)

	go accrualcalculate.WorkerPool(ctx, memStorageInterface, 10)
	router := Router(ctx, newHandStruct)

	return router, ctrl, newHandStruct, memStorageInterface
}
