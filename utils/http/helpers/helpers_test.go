package helpers_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var _ = Describe("helpers", func() {
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
