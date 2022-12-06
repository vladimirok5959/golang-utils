package pagination_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vladimirok5959/golang-utils/utils/pagination"
)

var _ = Describe("pagination", func() {
	Context("Data", func() {
		Context("New", func() {
			It("generate correct data", func() {
				pd := pagination.New(0, 0, 2)
				Expect(pd.CurrentPage()).To(Equal(int64(1)))
				Expect(pd.MaxPages()).To(Equal(int64(1)))
				Expect(pd.ResultsCount()).To(Equal(int64(0)))
				Expect(pd.ResultsPerPage()).To(Equal(int64(2)))
			})
		})

		Context("CurrentPage", func() {
			It("returns correct value", func() {
				Expect(pagination.New(0, 0, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(1, 0, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(2, 0, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(3, 0, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(-1, 0, 2).CurrentPage()).To(Equal(int64(1)))

				Expect(pagination.New(0, 2, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(1, 2, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(2, 2, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(3, 2, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(-1, 2, 2).CurrentPage()).To(Equal(int64(1)))

				Expect(pagination.New(0, 4, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(1, 4, 2).CurrentPage()).To(Equal(int64(1)))
				Expect(pagination.New(2, 4, 2).CurrentPage()).To(Equal(int64(2)))
				Expect(pagination.New(3, 4, 2).CurrentPage()).To(Equal(int64(2)))
				Expect(pagination.New(-1, 4, 2).CurrentPage()).To(Equal(int64(1)))
			})
		})

		Context("MaxPages", func() {
			It("returns correct value", func() {
				Expect(pagination.New(0, 0, 2).MaxPages()).To(Equal(int64(1)))
				Expect(pagination.New(0, 1, 2).MaxPages()).To(Equal(int64(1)))
				Expect(pagination.New(0, 2, 2).MaxPages()).To(Equal(int64(1)))
				Expect(pagination.New(0, 3, 2).MaxPages()).To(Equal(int64(2)))
				Expect(pagination.New(0, 4, 2).MaxPages()).To(Equal(int64(2)))
				Expect(pagination.New(0, 5, 2).MaxPages()).To(Equal(int64(3)))
				Expect(pagination.New(0, 6, 2).MaxPages()).To(Equal(int64(3)))
				Expect(pagination.New(0, 7, 2).MaxPages()).To(Equal(int64(4)))
				Expect(pagination.New(0, 8, 2).MaxPages()).To(Equal(int64(4)))
				Expect(pagination.New(0, 9, 2).MaxPages()).To(Equal(int64(5)))
				Expect(pagination.New(0, 10, 2).MaxPages()).To(Equal(int64(5)))
			})
		})

		Context("ResultsCount", func() {
			It("returns correct value", func() {
				Expect(pagination.New(0, 1, 0).ResultsCount()).To(Equal(int64(1)))
				Expect(pagination.New(0, 2, 0).ResultsCount()).To(Equal(int64(2)))
				Expect(pagination.New(0, 3, 0).ResultsCount()).To(Equal(int64(3)))
			})
		})

		Context("ResultsPerPage", func() {
			It("returns correct value", func() {
				Expect(pagination.New(0, 0, 1).ResultsPerPage()).To(Equal(int64(1)))
				Expect(pagination.New(0, 0, 2).ResultsPerPage()).To(Equal(int64(2)))
				Expect(pagination.New(0, 0, 3).ResultsPerPage()).To(Equal(int64(3)))
			})
		})

		Context("Limit", func() {
			It("returns correct value", func() {
				limit, offset := pagination.New(1, 10, 2).Limit()
				Expect(limit).To(Equal(int64(2)))
				Expect(offset).To(Equal(int64(0)))

				limit, offset = pagination.New(2, 10, 2).Limit()
				Expect(limit).To(Equal(int64(2)))
				Expect(offset).To(Equal(int64(2)))

				limit, offset = pagination.New(3, 10, 2).Limit()
				Expect(limit).To(Equal(int64(2)))
				Expect(offset).To(Equal(int64(4)))

				limit, offset = pagination.New(4, 10, 2).Limit()
				Expect(limit).To(Equal(int64(2)))
				Expect(offset).To(Equal(int64(6)))

				limit, offset = pagination.New(5, 10, 2).Limit()
				Expect(limit).To(Equal(int64(2)))
				Expect(offset).To(Equal(int64(8)))
			})
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pagination")
}
