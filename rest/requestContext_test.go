package rest_test

import (
	"github.com/gotgo/gokn/rest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestContext", func() {

	It("should add get remove", func() {
		ctx := rest.NewRequestContext()

		ctx.Add("a", "data", 37)
		result, found := ctx.Get("a", "data")

		Expect(found).To(BeTrue())
		Expect(result).To(Equal(37))

		result, found = ctx.Get("b", "data")
		Expect(found).To(BeFalse())

		ctx.Remove("a", "data")

		result, found = ctx.Get("b", "data")
		Expect(found).To(BeFalse())

	})
})
