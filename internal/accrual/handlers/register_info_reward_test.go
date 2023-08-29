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

func TestRegisterInfoReward_BadMethod(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/goods",
		method:     http.MethodGet,
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

func TestRegisterInfoReward_WrongContentType(t *testing.T) {
	data := struct {
		request    string
		method     string
		statusCode int
	}{
		request:    "/api/goods",
		method:     http.MethodPost,
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

func TestRegisterInfoReward_InvalidJSON(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       string
		statusCode int
	}{
		request:    "/api/goods",
		method:     http.MethodPost,
		body:       `{match: str`,
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

func TestRegisterInfoReward_RegisterInDBConflict(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       *models.Reward
		statusCode int
	}{
		request: "/api/goods",
		method:  http.MethodPost,
		body: &models.Reward{
			Match:      "Bork",
			Reward:     models.RewardDefault,
			RewardType: models.RewardTypeDefault,
		},
		statusCode: http.StatusConflict,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.TODO()

	pgErr := &pgconn.PgError{}
	pgErr.Code = pkg.UniqueViolationCode

	jsonBody, err := json.Marshal(data.body)
	require.NoError(t, err)

	storeInterface.EXPECT().RegisterInfoInDB(ctx, data.body).Return(pgErr)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(jsonBody), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterInfoReward_RegisterInDBFail(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       *models.Reward
		statusCode int
	}{
		request: "/api/goods",
		method:  http.MethodPost,
		body: &models.Reward{
			Match:      "Bork",
			Reward:     models.RewardDefault,
			RewardType: models.RewardTypeDefault,
		},
		statusCode: http.StatusConflict,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.TODO()

	pgErr := &pgconn.PgError{}

	jsonBody, err := json.Marshal(data.body)
	require.NoError(t, err)

	storeInterface.EXPECT().RegisterInfoInDB(ctx, data.body).Return(pgErr)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(jsonBody), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}

func TestRegisterInfoReward_PositiveCase(t *testing.T) {
	data := struct {
		request    string
		method     string
		body       *models.Reward
		statusCode int
	}{
		request: "/api/goods",
		method:  http.MethodPost,
		body: &models.Reward{
			Match:      "Bork",
			Reward:     models.RewardDefault,
			RewardType: models.RewardTypeDefault,
		},
		statusCode: http.StatusAccepted,
	}

	r, ctrl, _, storeInterface := runTestServer(t)
	defer ctrl.Finish()

	ts := httptest.NewServer(r)
	defer ts.Close()

	ctx := context.TODO()

	jsonBody, err := json.Marshal(data.body)
	require.NoError(t, err)

	storeInterface.EXPECT().RegisterInfoInDB(ctx, data.body).Return(nil)

	resp := testRequest(t, ts, data.method, data.request, bytes.NewReader(jsonBody), "", "application/json")
	defer resp.Body.Close()

	assert.Equal(t, data.statusCode, resp.StatusCode)
}
