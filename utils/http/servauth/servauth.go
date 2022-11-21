package servauth

import (
	"log"
	"net/http"
	"time"

	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var mRequests = &Requests{
	counter:   map[string]int{},
	lastTime:  map[string]int64{},
	cleanTime: time.Now().UTC().Unix(),
}

func BasicAuth(handler http.Handler, username, password, realm string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username != "" {
			mRequests.CleanupHourly()

			ip := helpers.ClientIP(r)
			reqs := mRequests.Count(ip)
			ltime := mRequests.Time(ip)

			// Reset counter
			if (time.Now().UTC().Unix() - ltime) >= 30 {
				reqs = 0
				mRequests.SetCount(ip, reqs)
				mRequests.SetTime(ip, time.Now().UTC().Unix())
			}

			// Restrict access
			if reqs >= 5 {
				w.Header().Set("Retry-After", "30")
				w.WriteHeader(429)
				if _, err := w.Write([]byte("Too Many Requests\n")); err != nil {
					log.Printf("%s\n", err.Error())
				}
				return
			}

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
				// Inc counter
				reqs = reqs + 1
				mRequests.SetCount(ip, reqs)
				mRequests.SetTime(ip, time.Now().UTC().Unix())

				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				if _, err := w.Write([]byte("Unauthorised\n")); err != nil {
					log.Printf("%s\n", err.Error())
				}
				return
			}

			if p != password {
				// Inc counter
				reqs = reqs + 1
				mRequests.SetCount(ip, reqs)
				mRequests.SetTime(ip, time.Now().UTC().Unix())

				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				if _, err := w.Write([]byte("Unauthorised\n")); err != nil {
					log.Printf("%s\n", err.Error())
				}
				return
			}

			// Reset counter
			if reqs > 0 {
				reqs = 0
				mRequests.SetCount(ip, reqs)
				mRequests.SetTime(ip, time.Now().UTC().Unix())
			}
		}

		handler.ServeHTTP(w, r)
	})
}
