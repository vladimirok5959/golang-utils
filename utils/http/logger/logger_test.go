package logger_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/logger"
)

var _ = Describe("logger", func() {
	Context("LogInternalError", func() {
		original := logger.ErrorLogFile
		logger.ErrorLogFile = "/tmp/test-err-out.log"

		f, err := os.OpenFile(logger.ErrorLogFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		Expect(err).To(Succeed())
		_, err = fmt.Fprint(f, "")
		Expect(err).To(Succeed())
		Expect(f.Close()).To(Succeed())

		logger.LogInternalError(fmt.Errorf("MyError 1"))
		file, err := os.ReadFile(logger.ErrorLogFile)
		Expect(err).To(Succeed())
		content := string(file)
		Expect(strings.HasSuffix(content, "MyError 1\n")).To(Equal(true))

		logger.LogInternalError(fmt.Errorf("MyError 2"))
		file, err = os.ReadFile(logger.ErrorLogFile)
		Expect(err).To(Succeed())
		content = string(file)
		Expect(strings.HasSuffix(content, "MyError 2\n")).To(Equal(true))

		// Restore original value
		logger.ErrorLogFile = original
	})

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
			Expect(buf.String()).To(ContainSubstring(`"GET /" 200`))
			Expect(buf.String()).To(ContainSubstring(`Go-http-client`))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "logger")
}
