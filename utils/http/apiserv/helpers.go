package apiserv

import "net/http"

// http.MethodDelete
// http.MethodGet
// http.MethodOptions
// http.MethodPatch
// http.MethodPost
// http.MethodPut

type TMethods []string

func Methods() TMethods {
	return []string{}
}

func (m TMethods) Delete() TMethods {
	return append(m, http.MethodDelete)
}

func (m TMethods) Get() TMethods {
	return append(m, http.MethodGet)
}

func (m TMethods) Options() TMethods {
	return append(m, http.MethodOptions)
}

func (m TMethods) Patch() TMethods {
	return append(m, http.MethodPatch)
}

func (m TMethods) Post() TMethods {
	return append(m, http.MethodPost)
}

func (m TMethods) Put() TMethods {
	return append(m, http.MethodPut)
}
