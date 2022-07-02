package logger

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (r *ResponseWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func ClientIP(r *http.Request) string {
	ips := ClientIPs(r)
	if len(ips) >= 1 {
		return ips[0]
	}
	return ""
}

func ClientIPs(r *http.Request) []string {
	ra := r.RemoteAddr
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), " "); xff != "" {
		ra = strings.Join([]string{xff, ra}, ",")
	}
	res := []string{}
	ips := strings.Split(ra, ",")
	for _, ip := range ips {
		res = append(res, strings.Trim(ip, " "))
	}
	return res
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
			ClientIP(r), r.Method, r.URL, nw.Status, time.Since(start).Seconds(),
		)
	})
}
