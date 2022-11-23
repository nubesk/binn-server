package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	w          http.ResponseWriter
	mw         io.Writer
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter, buf io.Writer) *loggingResponseWriter {
	return &loggingResponseWriter{
		w:          w,
		mw:         io.MultiWriter(w, buf),
		statusCode: http.StatusOK,
	}
}

func (lw *loggingResponseWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lw.mw.Write(b)
}

func (lw *loggingResponseWriter) WriteHeader(statusCode int) {
	lw.w.WriteHeader(statusCode)
	lw.statusCode = statusCode
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isErr := false
			reqBodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				isErr = true
			}
			r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))

			buf := bytes.NewBuffer([]byte{})
			lw := newLoggingResponseWriter(w, buf)
			next.ServeHTTP(lw, r)

			respBodyBytes, err := io.ReadAll(buf)
			if err != nil {
				isErr = true
			}

			if !isErr {
				logger.Printf("| %d | %#v | %#v \n", lw.statusCode, string(reqBodyBytes), string(respBodyBytes))
			}
		})
	}
}
