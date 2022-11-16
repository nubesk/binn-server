package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nubesk/binn"
)

func GetStreamHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		flusher, ok := w.(http.Flusher)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if logger != nil {
			logger.Printf("[out] connected")
		}
		closed := make(chan struct{}, 0)
		bn.Subscribe(func(b *binn.Bottle) bool {
			body := bottleToResponseItem(b)
			if bytes, err := json.Marshal(body); err == nil {
				if _, err := w.Write(bytes); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					closed <- struct{}{}
					return false
				}
				w.Write([]byte("\n"))
				flusher.Flush()
				if logger != nil {
					closed <- struct{}{}
					logger.Printf("[out] {\"message\": \"%s\"}\n", b.Msg)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return false
			}
			flusher.Flush()
			return true
		})
		select {
		case <-r.Context().Done():
		case <-closed:
		}
		if logger != nil {
			logger.Printf("[out] disconnected")
		}
		return
	}
}
