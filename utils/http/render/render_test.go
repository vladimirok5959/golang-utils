package render_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
	"github.com/vladimirok5959/golang-utils/utils/http/render"
)

var _ = Describe("render", func() {
	var w *helpers.FakeResponseWriter
	var r *http.Request

	BeforeEach(func() {
		w = helpers.NewFakeResponseWriter()
		r = &http.Request{}
	})

	Context("HTML", func() {
		It("render", func() {
			var data struct {
				URL string
			}
			data.URL = "/example.html"
			Expect(render.HTML(w, r, nil, &data, "Url: {{$.Data.URL}}", http.StatusNotFound)).To(BeTrue())
			Expect(w.Body).To(Equal([]byte("Url: /example.html")))
			Expect(w.Headers).To(Equal(http.Header{
				"Content-Type": []string{"text/html"},
			}))
			Expect(w.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("JSON", func() {
		It("render", func() {
			var data struct{ Field string }
			Expect(render.JSON(w, r, data)).To(BeTrue())
			Expect(w.Body).To(Equal([]byte(`{"Field":""}`)))
			Expect(w.Headers).To(Equal(http.Header{
				"Content-Type": []string{"application/json"},
			}))
			Expect(w.StatusCode).To(Equal(http.StatusOK))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "render")
}
