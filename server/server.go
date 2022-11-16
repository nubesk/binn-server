package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nubesk/binn"
	"github.com/nubesk/binn-server/server/handler"
)

func New(bn *binn.Binn, addr string, logger *log.Logger) *http.Server {
	r := chi.NewRouter()
	r.Get("/", handler.GetHandlerFunc(bn, logger))
	r.Post("/", handler.PostHandlerFunc(bn, logger))
	r.Get("/stream", handler.GetStreamHandlerFunc(bn, logger))

	rr := chi.NewRouter()
	rr.Mount("/api/bottles", r)

	return &http.Server{
		Addr:    addr,
		Handler: rr,
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
