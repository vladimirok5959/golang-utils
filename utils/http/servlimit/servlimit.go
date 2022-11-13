package servlimit

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

func ReqPerSecond(handler http.Handler, requests int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cleanup every hour
		if (time.Now().UTC().Unix() - mRequests.cleanTime) > 3600 {
			mRequests.Cleanup()
		}

		ip := helpers.ClientIP(r)
		reqs := mRequests.Count(ip)
		ltime := mRequests.Time(ip)

		// Inc counter
		reqs = reqs + 1
		mRequests.SetCount(ip, reqs)

		// Reset counter
		if (time.Now().UTC().Unix() - ltime) >= 1 {
			reqs = 0
			mRequests.SetCount(ip, 0)
		}

		// Restrict access
		if reqs >= requests {
			w.WriteHeader(429)
			if _, err := w.Write([]byte("Too Many Requests\n")); err != nil {
				log.Printf("%s\n", err.Error())
			}
			return
		}

		mRequests.SetTime(ip, time.Now().UTC().Unix())

		handler.ServeHTTP(w, r)
	})
}
