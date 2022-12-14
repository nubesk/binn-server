package handler

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nubesk/binn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStreamHandler(t *testing.T) {
	storage := binn.NewBottleStorage(2)
	err := storage.Add(&binn.Bottle{
		Msg: "sample message",
	})
	require.NoError(t, err)
	bn := binn.New(storage, 0)
	handler := GetStreamHandlerFunc(bn, nil)
	req := httptest.NewRequest("GET", "http://example.com/api/bottles", nil)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	bs := strings.Split(string(body), "\n")
	require.GreaterOrEqual(t, len(bs), 1)
	assert.Equal(t, `{"message":"sample message"}`, bs[0])
}
