package logger

import (
	"net/http"
	"strings"
)

var RollBarEnabled = false

type RollBarStatusCodes []int

var RollBarSkipStatusCodes = RollBarStatusCodes{
	http.StatusContinue,           // 100
	http.StatusSwitchingProtocols, // 101
	http.StatusProcessing,         // 102
	http.StatusEarlyHints,         // 103

	http.StatusOK,                   // 200
	http.StatusCreated,              // 201
	http.StatusAccepted,             // 202
	http.StatusNonAuthoritativeInfo, // 203
	http.StatusNoContent,            // 204
	http.StatusResetContent,         // 205
	http.StatusPartialContent,       // 206
	http.StatusMultiStatus,          // 207
	http.StatusAlreadyReported,      // 208
	http.StatusIMUsed,               // 226

	http.StatusMultipleChoices,   // 300
	http.StatusMovedPermanently,  // 301
	http.StatusFound,             // 302
	http.StatusSeeOther,          // 303
	http.StatusNotModified,       // 304
	http.StatusUseProxy,          // 305
	http.StatusTemporaryRedirect, // 307
	http.StatusPermanentRedirect, // 308

	http.StatusForbidden,        // 403
	http.StatusNotFound,         // 404
	http.StatusMethodNotAllowed, // 405
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
