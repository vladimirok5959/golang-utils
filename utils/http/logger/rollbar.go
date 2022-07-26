package logger

import (
	"net/http"
	"strings"
)

var RollBarEnabled = false

type RollBarStatusCodes []int

var RollBarSkipStatusCodes = RollBarStatusCodes{
	http.StatusOK,
	http.StatusNotModified,
	http.StatusTemporaryRedirect,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
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
