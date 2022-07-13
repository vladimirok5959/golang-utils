package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

// func ClientIP(r *http.Request) string
// func ClientIPs(r *http.Request) []string
// func HandleAppStatus() http.Handler
// func HandleFile(data, contentType string) http.Handler
// func HandleImagePng(data string) http.Handler
// func HandleTextCss(data string) http.Handler
// func HandleTextJavaScript(data string) http.Handler
// func HandleTextPlain(data string) http.Handler
// func HandleTextXml(data string) http.Handler
// func MinifyHtmlCode(str string) string
// func RespondAsBadRequest(w http.ResponseWriter, r *http.Request, err error)
// func RespondAsMethodNotAllowed(w http.ResponseWriter, r *http.Request)
// func SessionStart(w http.ResponseWriter, r *http.Request) (*session.Session, error)
// func SetLanguageCookie(w http.ResponseWriter, r *http.Request) error

var mHtml = regexp.MustCompile(`>[\n\t\r]+<`)
var mHtmlLeft = regexp.MustCompile(`>[\n\t\r]+`)
var mHtmlRight = regexp.MustCompile(`[\n\t\r]+<`)

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

func HandleAppStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RespondAsMethodNotAllowed(w, r)
			return
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		type respMemory struct {
			Alloc      uint64 `json:"alloc"`
			NumGC      uint32 `json:"num_gc"`
			Sys        uint64 `json:"sys"`
			TotalAlloc uint64 `json:"total_alloc"`
		}

		type respRoot struct {
			Memory   respMemory `json:"memory"`
			Routines int        `json:"routines"`
		}

		resp := respRoot{
			Memory: respMemory{
				Alloc:      m.Alloc,
				NumGC:      m.NumGC,
				Sys:        m.Sys,
				TotalAlloc: m.TotalAlloc,
			},
			Routines: runtime.NumGoroutine(),
		}

		j, err := json.Marshal(resp)
		if err != nil {
			RespondAsBadRequest(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(j); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	})
}

func HandleFile(data, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RespondAsMethodNotAllowed(w, r)
			return
		}
		w.Header().Set("Content-Type", contentType)
		if _, err := w.Write([]byte(data)); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	})
}

func HandleImagePng(data string) http.Handler {
	return HandleFile(data, "image/png")
}

func HandleTextCss(data string) http.Handler {
	return HandleFile(data, "text/css")
}

func HandleTextJavaScript(data string) http.Handler {
	return HandleFile(data, "text/javascript")
}

func HandleTextPlain(data string) http.Handler {
	return HandleFile(data, "text/plain")
}

func HandleTextXml(data string) http.Handler {
	return HandleFile(data, "text/xml")
}

func MinifyHtmlCode(str string) string {
	str = mHtml.ReplaceAllString(str, "><")
	str = mHtmlLeft.ReplaceAllString(str, ">")
	str = mHtmlRight.ReplaceAllString(str, "<")
	return str
}

func RespondAsBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		log.Printf("%s\n", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if _, e := w.Write([]byte(`{"error":` + strconv.Quote(err.Error()) + `}`)); e != nil {
			log.Printf("%s\n", e.Error())
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func RespondAsMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Example:
//
// sess, err := helpers.SessionStart(w, r)
//
// if err != nil && !errors.Is(err, os.ErrNotExist) {
//
//	helpers.RespondAsBadRequest(w, r, err)
// 	return
//
// }
//
// defer sess.Close()
func SessionStart(w http.ResponseWriter, r *http.Request) (*session.Session, error) {
	sess, err := session.New(w, r, "/tmp")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("%s\n", err.Error())
	}
	return sess, err
}

// Example:
//
// if err = r.ParseForm(); err != nil {
// 	helpers.RespondAsBadRequest(w, r, err)
// 	return
// }
//
// if err = helpers.SetLanguageCookie(w, r); err != nil {
// 		helpers.RespondAsBadRequest(w, r, err)
// 		return
// }
func SetLanguageCookie(w http.ResponseWriter, r *http.Request) error {
	var clang string
	if c, err := r.Cookie("lang"); err == nil {
		clang = c.Value
	}
	lang := r.Form.Get("lang")
	if lang != "" && lang != clang {
		http.SetCookie(w, &http.Cookie{
			Name:     "lang",
			Value:    lang,
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		})
	}
	return nil
}
