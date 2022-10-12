package servtools

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/vladimirok5959/golang-utils/utils/http/render"
	"github.com/vladimirok5959/golang-utils/utils/penv"
)

var (
	//go:embed config-vars.html
	configVarsHtml string
)

// config - must be a pointer to config structure
func ConfigVars(config any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !render.HTML(
			w,
			r,
			template.FuncMap{
				"secret": func(value string) template.HTML {
					return template.HTML(`<span onclick="if(this.innerHTML=='**********'){this.innerHTML='` + value + `';}else{this.innerHTML='**********';}" style="cursor:pointer;">**********</span>`)
				},
			},
			penv.DumpConfig(config),
			configVarsHtml,
			http.StatusOK,
		) {
			return
		}
	})
}
