package logger

import (
	"log"
	"net/http"
	"time"

	"github.com/rollbar/rollbar-go"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var RollBarEnabled = false

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
		if RollBarEnabled {
			if !(nw.Status == http.StatusOK || nw.Status == http.StatusNotFound) {
				rollbar.Error(r, nw)
			}
		}
		log.Printf(
			"\"%s\" \"%s %s\" %d \"%.3f ms\"\n",
			helpers.ClientIP(r), r.Method, r.URL, nw.Status, time.Since(start).Seconds(),
		)
	})
}
