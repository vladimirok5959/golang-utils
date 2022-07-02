package helpers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"

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

func SessionStart(w http.ResponseWriter, r *http.Request) (*session.Session, error) {
	sess, err := session.New(w, r, "/tmp")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("%s\n", err.Error())
	}
	return sess, err
}
