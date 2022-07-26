package logger

import "net/http"

var RollBarEnabled = false

type RollBarStatusCodes []int

var RollBarSkipStatusCodes = RollBarStatusCodes{
	http.StatusOK,
	http.StatusNotModified,
	http.StatusTemporaryRedirect,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
}

func (r RollBarStatusCodes) Contain(status int) bool {
	for _, v := range r {
		if v == status {
			return true
		}
	}
	return false
}
