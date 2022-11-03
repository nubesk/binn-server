package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nubesk/binn"
)

type requestBody struct {
	Msg string `json:"message"`
}

type responseBody struct {
	Msg string `json:"message"`
}

func New(bn *binn.Binn, addr string, logger *log.Logger) *http.Server {
	r := chi.NewRouter()
	r.Get("/", bottleGetHandlerFunc(bn, logger))
	r.Post("/", bottlePostHandlerFunc(bn, logger))

	rr := chi.NewRouter()
	rr.Mount("/api/bottles", r)

	return &http.Server{
		Addr:    addr,
		Handler: rr,
	}
}

func bottleToResponse(b *binn.Bottle) *responseBody {
	return &responseBody{
		Msg: b.Msg,
	}
}

func requestToBottle(body *requestBody) *binn.Bottle {
	return &binn.Bottle{
		Msg: body.Msg,
	}
}

func bottleGetHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
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
		ch := make(chan *binn.Bottle, 1)
		err := bn.Subscribe(ch)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	Loop:
		for {
			select {
			case <-r.Context().Done():
				break Loop
			case b, ok := <-ch:
				if !ok {
					ch := make(chan *binn.Bottle, 1)
					err := bn.Subscribe(ch)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				body := bottleToResponse(b)
				if bytes, err := json.Marshal(body); err == nil {
					if _, err := w.Write(bytes); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.Write([]byte("\n"))
					flusher.Flush()
					if logger != nil {
						logger.Printf("[out] {\"message\": \"%s\"}\n", b.Msg)
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				flusher.Flush()
			}
		}
		if logger != nil {
			logger.Printf("[out] disconnected")
		}
		return
	}
}

func bottlePostHandlerFunc(bn *binn.Binn, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var body requestBody
		if err := json.Unmarshal(bytes, &body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b := requestToBottle(&body)
		err = bn.Publish(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if bytes, err := json.Marshal(body); err == nil {
			if _, err := w.Write(bytes); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if logger != nil {
				logger.Printf("[in] {\"message\": \"%s\"}\n", b.Msg)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "GET, POST")
		next.ServeHTTP(w, r)
	})
}
