package apiserv

import "net/http"

type Handler struct {
	handler http.Handler
	methods []string
}

func (s Handler) IsMethod(method string) bool {
	for _, v := range s.methods {
		if v == method {
			return true
		}
	}
	return false
}
