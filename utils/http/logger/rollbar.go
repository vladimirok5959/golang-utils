package logger

import (
	"net/http"
	"strings"
)

var RollBarEnabled = false

type RollBarStatusCodes []int

var RollBarSkipStatusCodes = RollBarStatusCodes{
	http.StatusForbidden,
	http.StatusMethodNotAllowed,
	http.StatusNotFound,
	http.StatusNotModified,
	http.StatusOK,
	http.StatusTemporaryRedirect,
}

func (r RollBarStatusCodes) contain(status int) bool {
	for _, v := range r {
		if v == status {
			return true
		}
	}
	return false
}

type RollBarErrors []string

var RollBarSkipErrors = RollBarErrors{}

func (r RollBarErrors) contain(str string) bool {
	for _, v := range r {
		if v == str || strings.Contains(str, v) {
			return true
		}
	}
	return false
}
