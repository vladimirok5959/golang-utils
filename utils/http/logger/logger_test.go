package logger_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/logger"
)

var _ = Describe("logger", func() {
	Context("LogRequests", func() {
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
			srv = httptest.NewServer(logger.LogRequests(getTestHandler()))
			client = srv.Client()
		})

		AfterEach(func() {
			srv.Close()
		})

		It("log http requests", func() {
			buf := new(bytes.Buffer)
			log.SetOutput(buf)

			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(buf.String()).To(ContainSubstring("\"GET /\" 200"))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "logger")
}
