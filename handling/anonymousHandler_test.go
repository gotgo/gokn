package handling_test

import (
	. "github.com/gotgo/gokn/handling"
	"github.com/gotgo/gokn/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AnonymousHandler", func() {

	It("should work", func() {
		count := 0
		handler := func(*rest.Request, rest.Responder) {
			count++
		}
		wrapper := AnonymousHandler(handler)
		wrapper(nil, nil)
		Expect(count).To(Equal(1))
	})

})
