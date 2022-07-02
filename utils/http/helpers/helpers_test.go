package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/http/helpers"
)

var _ = Describe("helpers", func() {
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
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "helpers")
}
