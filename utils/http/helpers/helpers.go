package helpers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/vladimirok5959/golang-server-sessions/session"
)

var mHtml = regexp.MustCompile(`>[\n\t\r\s]+<`)
var mHtmlLeft = regexp.MustCompile(`>[\n\t\r\s]+`)
var mHtmlRight = regexp.MustCompile(`[\n\t\r\s]+<`)

func MinifyHtmlCode(str string) string {
	str = mHtml.ReplaceAllString(str, "><")
	str = mHtmlLeft.ReplaceAllString(str, ">")
	str = mHtmlRight.ReplaceAllString(str, "<")
	return str
}

// Example:
//
// sess := helpers.SessionStart(w, r)
//
// defer sess.Close()
func SessionStart(w http.ResponseWriter, r *http.Request) (*session.Session, error) {
	sess, err := session.New(w, r, "/tmp")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("%s\n", err.Error())
	}
	return sess, err
}

func RespondAsBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		log.Printf("%s\n", err.Error())
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if _, e := w.Write([]byte(`{"error":"` + strconv.Quote(err.Error()) + `"}`)); e != nil {
			log.Printf("%s\n", e.Error())
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func RespondAsMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
