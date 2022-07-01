package utils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils"
)

var _ = Describe("utils", func() {
	Context("RandomString", func() {
		It("generate random string with defined length", func() {
			Expect(len(utils.RandomString(10))).To(Equal(10))
			Expect(len(utils.RandomString(20))).To(Equal(20))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils")
}
