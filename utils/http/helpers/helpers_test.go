package helpers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var _ = Describe("helpers", func() {
	var srv *httptest.Server
	var client *http.Client

	Context("ClientIP", func() {
		It("return client IP", func() {
			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "127.0.0.1",
			})).To(Equal("127.0.0.1"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.1,127.0.0.1",
			})).To(Equal("192.168.0.1"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.1, 127.0.0.1",
			})).To(Equal("192.168.0.1"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.50,192.168.0.1,127.0.0.1",
			})).To(Equal("192.168.0.50"))

			Expect(helpers.ClientIP(&http.Request{
				RemoteAddr: "192.168.0.50, 192.168.0.1, 127.0.0.1",
			})).To(Equal("192.168.0.50"))
		})
	})

	Context("ClientIPs", func() {
		It("return array of client IPs", func() {
			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "127.0.0.1",
			})).To(ConsistOf(
				"127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.1,127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.1", "127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.1, 127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.1", "127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.50,192.168.0.1,127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.50", "192.168.0.1", "127.0.0.1",
			))

			Expect(helpers.ClientIPs(&http.Request{
				RemoteAddr: "192.168.0.50, 192.168.0.1, 127.0.0.1",
			})).To(ConsistOf(
				"192.168.0.50", "192.168.0.1", "127.0.0.1",
			))
		})
	})

	Context("HandleAppStatus", func() {
		BeforeEach(func() {
			srv = httptest.NewServer(helpers.HandleAppStatus())
			client = srv.Client()
		})

		AfterEach(func() {
			srv.Close()
		})

		It("handle app status", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

			body, err := io.ReadAll(resp.Body)
			Expect(err).To(Succeed())

			Expect(string(body)).To(MatchRegexp(`{"memory":{"alloc":[0-9]+,"num_gc":[0-9]+,"sys":[0-9]+,"total_alloc":[0-9]+},"routines":[0-9]+}`))
		})
	})

	Context("MinifyHtmlCode", func() {
		It("minify Html code", func() {
			Expect(helpers.MinifyHtmlCode(`
				<!doctype html>
				<html lang="uk">
					<head>
						<meta charset="utf-8" />
					</head>
					<body>
						Index
					</body>
				</html>
			`)).To(Equal(`<!doctype html><html lang="uk"><head><meta charset="utf-8" /></head><body>Index</body></html>`))

			Expect(helpers.MinifyHtmlCode(`
				<div>
					<a href="#">Link 1</a>, <a href="#">Link 2</a>
				</div>
			`)).To(Equal(`<div><a href="#">Link 1</a>, <a href="#">Link 2</a></div>`))

			Expect(helpers.MinifyHtmlCode(`
				<div>
					<b>Contacts:</b> <a href="#">Link 1</a>, <a href="#">Link 2</a>
				</div>
			`)).To(Equal(`<div><b>Contacts:</b> <a href="#">Link 1</a>, <a href="#">Link 2</a></div>`))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "helpers")
}
