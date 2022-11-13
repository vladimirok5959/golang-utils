package servlimit_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/servlimit"
)

var _ = Describe("servlimit", func() {
	Context("ReqPerSecond", func() {
		var srv *httptest.Server
		var client *http.Client

		var getTestHandler = func() http.HandlerFunc {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte("Index")); err != nil {
					fmt.Printf("%s\n", err.Error())
				}
			})
		}

		BeforeEach(func() {
			servlimit.MRequests.Cleanup()
			srv = httptest.NewServer(servlimit.ReqPerSecond(getTestHandler(), 1))
			client = srv.Client()
		})

		AfterEach(func() {
			srv.Close()
		})

		It("process request", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Index"))
		})

		It("process multiple requests", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			time.Sleep(1 * time.Second)

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Index"))
		})

		It("block multiple requests", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTooManyRequests))
			Expect(resp.Header.Get("Retry-After")).To(Equal("1"))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Too Many Requests\n"))
		})

		It("block more multiple requests", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTooManyRequests))
			Expect(resp.Header.Get("Retry-After")).To(Equal("1"))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTooManyRequests))
			Expect(resp.Header.Get("Retry-After")).To(Equal("1"))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTooManyRequests))
			Expect(resp.Header.Get("Retry-After")).To(Equal("1"))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTooManyRequests))
			Expect(resp.Header.Get("Retry-After")).To(Equal("1"))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Too Many Requests\n"))
		})

		It("clean requests data in memory", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			servlimit.MRequests.Cleanup()

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Index"))
		})

		It("process 3 requests per second", func() {
			srv.Close()
			srv = httptest.NewServer(servlimit.ReqPerSecond(getTestHandler(), 3))
			client = srv.Client()

			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Retry-After")).To(Equal(""))

			resp, err = client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTooManyRequests))
			Expect(resp.Header.Get("Retry-After")).To(Equal("1"))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Too Many Requests\n"))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "servlimit")
}
