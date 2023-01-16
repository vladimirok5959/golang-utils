package apiserv_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/apiserv"
)

var _ = Describe("apiserv", func() {
	Context("Methods", func() {
		It("build correct array", func() {
			m := apiserv.Methods()
			Expect(m).To(Equal(apiserv.TMethods{}))

			// Single
			Expect(m.Delete()).To(Equal(apiserv.TMethods{http.MethodDelete}))
			Expect(m.Get()).To(Equal(apiserv.TMethods{http.MethodGet}))
			Expect(m.Options()).To(Equal(apiserv.TMethods{http.MethodOptions}))
			Expect(m.Patch()).To(Equal(apiserv.TMethods{http.MethodPatch}))
			Expect(m.Post()).To(Equal(apiserv.TMethods{http.MethodPost}))
			Expect(m.Put()).To(Equal(apiserv.TMethods{http.MethodPut}))

			// Multiple
			Expect(m.Get().Put()).To(Equal(apiserv.TMethods{
				http.MethodGet,
				http.MethodPut,
			}))

			Expect(m.Get().Post()).To(Equal(apiserv.TMethods{
				http.MethodGet,
				http.MethodPost,
			}))

			Expect(m.Get().Put().Post().Delete()).To(Equal(apiserv.TMethods{
				http.MethodGet,
				http.MethodPut,
				http.MethodPost,
				http.MethodDelete,
			}))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "apiserv")
}
