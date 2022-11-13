package servauth

import (
	"log"
	"net/http"
)

// TODO: protect from bruteforce

func BasicAuth(handler http.Handler, username, password, realm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username != "" {
			if realm == "" {
				realm = "Please enter username and password"
			}

			u, p, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				if _, err := w.Write([]byte("Unauthorised\n")); err != nil {
					log.Printf("%s\n", err.Error())
				}
				return
			}

			if u != username {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				if _, err := w.Write([]byte("Unauthorised\n")); err != nil {
					log.Printf("%s\n", err.Error())
				}
				return
			}

			if p != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				if _, err := w.Write([]byte("Unauthorised\n")); err != nil {
					log.Printf("%s\n", err.Error())
				}
				return
			}
		}

		handler.ServeHTTP(w, r)
	})
}
