package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rollbar/rollbar-go"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var AccessLogFile = ""
var ErrorLogFile = ""

func appendToLogFile(fileName, msg string) error {
	flags := os.O_RDWR | os.O_CREATE | os.O_APPEND
	f, err := os.OpenFile(fileName, flags, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := fmt.Fprint(f, msg); err != nil {
		return err
	}
	return nil
}

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
		msg := fmt.Sprintf(
			"\"%s\" \"%s %s\" %d %d \"%.3f ms\" \"%s\"\n",
			strings.Join(helpers.ClientIPs(r), ", "),
			r.Method,
			r.URL,
			nw.Status,
			nw.Size,
			time.Since(start).Seconds(),
			ua,
		)
		log.Printf("%s", msg)

		if nw.Status < 400 {
			if AccessLogFile != "" {
				if err := appendToLogFile(AccessLogFile, start.Format("2009/01/23 01:23:23")+" "+msg); err != nil {
					log.Printf("%s\n", err.Error())
				}
			}
		} else {
			if ErrorLogFile != "" {
				if err := appendToLogFile(ErrorLogFile, start.Format("2009/01/23 01:23:23")+" "+msg); err != nil {
					log.Printf("%s\n", err.Error())
				}
			}
		}

		if RollBarEnabled && !RollBarSkipStatusCodes.contain(nw.Status) {
			if !RollBarSkipErrors.contain(string(nw.Content)) {
				rollbar.Error(r, nw.Status, nw.Size, string(nw.Content))
			}
		}
	})
}
