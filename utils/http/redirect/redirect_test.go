package redirect_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/redirect"
)

var _ = Describe("redirect", func() {
	Context("Handler", func() {
		var srv *httptest.Server
		var client *http.Client

		BeforeEach(func() {
			srv = httptest.NewServer(redirect.Handler("/example.html"))
			client = srv.Client()

			// Disable HTTP redirects
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		})

		AfterEach(func() {
			srv.Close()
		})

		It("doing redirect", func() {
			resp, err := client.Get(srv.URL + "/")
			Expect(err).To(Succeed())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusTemporaryRedirect))
			Expect(resp.Header.Get("Location")).To(Equal("/example.html"))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "redirect")
}
