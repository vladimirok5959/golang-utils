package render

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
	"github.com/vladimirok5959/golang-utils/utils/http/logger"
)

// func HTML(w http.ResponseWriter, r *http.Request, f template.FuncMap, d interface{}, s string, httpStatusCode int) bool
// func JSON(w http.ResponseWriter, r *http.Request, o interface{}) bool

func HTML(w http.ResponseWriter, r *http.Request, f template.FuncMap, d interface{}, s string, httpStatusCode int) bool {
	tmpl := template.New("tmpl")
	if f != nil {
		tmpl = tmpl.Funcs(f)
	}
	var err error
	tmpl, err = tmpl.Parse(s)
	if err != nil {
		helpers.RespondAsBadRequest(w, r, err)
		return false
	}
	type Response struct {
		Data interface{}
	}
	var html bytes.Buffer
	if err := tmpl.Execute(&html, Response{Data: d}); err != nil {
		helpers.RespondAsBadRequest(w, r, err)
		return false
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(httpStatusCode)
	if _, err := w.Write([]byte(helpers.MinifyHtmlCode(html.String()))); err != nil {
		logger.LogInternalError(err)
	}
	return true
}

func JSON(w http.ResponseWriter, r *http.Request, o interface{}) bool {
	j, err := json.Marshal(o)
	if err != nil {
		helpers.RespondAsBadRequest(w, r, err)
		return false
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(j); err != nil {
		logger.LogInternalError(err)
	}
	return true
}
