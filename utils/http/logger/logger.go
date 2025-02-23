package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
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

func LogError(format string, a ...any) {
	msg := fmt.Sprintf("[ERROR] %s\n", fmt.Sprintf(format, a...))
	if pc, file, line, ok := runtime.Caller(1); ok {
		msg = fmt.Sprintf(
			"[ERROR] (%s, line #%d, func: %v) %s\n",
			file,
			line,
			runtime.FuncForPC(pc).Name(),
			fmt.Sprintf(format, a...),
		)
	}
	log.Printf("%s", msg)
	if ErrorLogFile != "" {
		if err := appendToLogFile(ErrorLogFile, time.Now().Format("2006/01/02 15:04:05")+" "+msg); err != nil {
			log.Printf("%s\n", err.Error())
		}
	}
}

func LogInfo(format string, a ...any) {
	msg := fmt.Sprintf("[INFO] %s\n", fmt.Sprintf(format, a...))
	log.Printf("%s", msg)
	if AccessLogFile != "" {
		if err := appendToLogFile(AccessLogFile, time.Now().Format("2006/01/02 15:04:05")+" "+msg); err != nil {
			log.Printf("%s\n", err.Error())
		}
	}
}

func LogInternalError(err error) {
	msg := fmt.Sprintf("[ERROR] %s\n", err.Error())
	if pc, file, line, ok := runtime.Caller(1); ok {
		msg = fmt.Sprintf(
			"[ERROR] (%s, line #%d, func: %v) %s\n",
			file,
			line,
			runtime.FuncForPC(pc).Name(),
			err.Error(),
		)
	}
	log.Printf("%s", msg)
	if ErrorLogFile != "" {
		if err := appendToLogFile(ErrorLogFile, time.Now().Format("2006/01/02 15:04:05")+" "+msg); err != nil {
			log.Printf("%s\n", err.Error())
		}
	}
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

		if nw.Status < 400 {
			log.Printf("[ACCESS] %s", msg)
			if AccessLogFile != "" {
				if err := appendToLogFile(AccessLogFile, start.Format("2006/01/02 15:04:05")+" "+msg); err != nil {
					log.Printf("%s\n", err.Error())
				}
			}
		} else {
			log.Printf("[ERROR] %s", msg)
			if ErrorLogFile != "" {
				if err := appendToLogFile(ErrorLogFile, start.Format("2006/01/02 15:04:05")+" "+msg); err != nil {
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
