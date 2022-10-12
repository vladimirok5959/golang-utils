package servauth_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/servauth"
)

var _ = Describe("servauth", func() {
	Context("BasicAuth", func() {
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
			srv = httptest.NewServer(servauth.BasicAuth(getTestHandler(), "user", "pass", "msg"))
			client = srv.Client()
		})

		AfterEach(func() {
			srv.Close()
		})

		It("request credentials", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			Expect(resp.Header["Www-Authenticate"]).To(Equal([]string{`Basic realm="msg"`}))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Unauthorised\n"))
		})

		It("show with correct credentials", func() {
			req, err := http.NewRequest("GET", srv.URL+"/", nil)
			Expect(err).To(Succeed())
			req.SetBasicAuth("user", "pass")

			resp, err := client.Do(req)
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Index"))
		})

		It("don't request credentials on empty username", func() {
			srv.Close()
			srv = httptest.NewServer(servauth.BasicAuth(getTestHandler(), "", "", ""))
			client = srv.Client()

			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Index"))
		})

		It("request credentials on not empty username but empty password", func() {
			srv.Close()
			srv = httptest.NewServer(servauth.BasicAuth(getTestHandler(), "user", "", "msg"))
			client = srv.Client()

			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			Expect(resp.Header["Www-Authenticate"]).To(Equal([]string{`Basic realm="msg"`}))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())
			Expect(string(body)).To(Equal("Unauthorised\n"))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "servauth")
}
