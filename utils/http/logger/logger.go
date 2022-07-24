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
	Content []byte
	Status  int
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if RollBarEnabled {
		if !(w.Status == http.StatusOK || w.Status == http.StatusNotFound) {
			w.Content = append(w.Content, b...)
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func LogRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nw := &ResponseWriter{
			ResponseWriter: w,
			Content:        []byte{},
			Status:         http.StatusOK,
		}
		handler.ServeHTTP(nw, r)
		if RollBarEnabled {
			if !(nw.Status == http.StatusOK || nw.Status == http.StatusNotFound) {
				rollbar.Error(r, string(nw.Content))
			}
		}
		log.Printf(
			"\"%s\" \"%s %s\" %d \"%.3f ms\"\n",
			helpers.ClientIP(r), r.Method, r.URL, nw.Status, time.Since(start).Seconds(),
		)
	})
}
