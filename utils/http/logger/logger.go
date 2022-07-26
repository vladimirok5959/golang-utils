package logger

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rollbar/rollbar-go"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

func LogInternalError(err error) {
	log.Printf("%s\n", err.Error())
	if RollBarEnabled && !RollBarSkipErrors.contain(err.Error()) {
		rollbar.Error(err)
	}
}

func LogRequests(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nw := &ResponseWriter{
			ResponseWriter: w,
			Content:        []byte{},
			Size:           0,
			Status:         http.StatusOK,
		}
		handler.ServeHTTP(nw, r)
		ua := strings.TrimSpace(r.Header.Get("User-Agent"))
		if ua == "" || len(ua) > 256 {
			ua = "-"
		}
		log.Printf(
			"\"%s\" \"%s %s\" %d %d \"%.3f ms\" \"%s\"\n",
			strings.Join(helpers.ClientIPs(r), ", "),
			r.Method,
			r.URL,
			nw.Status,
			nw.Size,
			time.Since(start).Seconds(),
			ua,
		)
		if RollBarEnabled && !RollBarSkipStatusCodes.contain(nw.Status) {
			if !RollBarSkipErrors.contain(string(nw.Content)) {
				rollbar.Error(r, nw.Status, nw.Size, string(nw.Content))
			}
		}
	})
}
