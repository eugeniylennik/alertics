package handlers_test

import (
	"context"
	"github.com/eugeniylennik/alertics/internal/database"
	"github.com/eugeniylennik/alertics/internal/router"
	"github.com/eugeniylennik/alertics/internal/server"
	"github.com/eugeniylennik/alertics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_RecordMetrics(t *testing.T) {
	cfg := server.InitConfigServer()
	m := storage.NewMemStorage(cfg.StoreFile, cfg.StoreInterval == 0)
	client, err := database.NewClient(context.TODO(), 5, cfg.Dsn)
	if err != nil {
		log.Fatalln(err)
	}
	r := router.NewRouter(m, client)
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/update/gauge/Alloc/12.12")
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, body, "")
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
