package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nubesk/binn"
)

func TestHandlePostBottle(t *testing.T) {
	storage := binn.NewBottleStorage(1)
	bn := binn.New(storage, 0)
	handler := PostHandlerFunc(bn, nil)
	reqBody := bytes.NewBufferString(`{"message":"sample message"}`)
	req := httptest.NewRequest("POST", "http://example.com/api/bottles", reqBody)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	assert.Equal(t, 200, resp.StatusCode)
	b, err := storage.Get()
	require.NoError(t, err)
	assert.Equal(t, "sample message", b.Msg)
}
