package logger

import (
	"log"
	"net/http"
	"time"

	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (r *ResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func LogRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nw := &ResponseWriter{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		handler.ServeHTTP(nw, r)
		log.Printf(
			"\"%s\" \"%s %s\" %d \"%.3f ms\"\n",
			helpers.ClientIP(r), r.Method, r.URL, nw.Status, time.Since(start).Seconds(),
		)
	})
}
