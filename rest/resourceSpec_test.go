package rest_test

import (
	"reflect"

	"github.com/gotgo/gokn/rest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestMessage struct {
	Message string `json:"message"`
}

var _ = Describe("ResourceSpec", func() {

	var (
		def  *rest.ResourceDef
		spec *rest.ResourceSpec
	)

	BeforeEach(func() {
		def = &rest.ResourceDef{
			ResourceT:    "/abc/{id}",
			ResourceArgs: nil,
			Verb:         "GET",
			Headers:      nil,
			RequestBody:  reflect.TypeOf(TestMessage{}),
		}

		spec = rest.NewResourceSpec("application/json")
		spec.Use(def)
	})

	Context("client behavior", func() {
		It("should handle a nil argument, and return the spec with it's template", func() {
			var cl rest.ClientResource
			cl = spec

			request := cl.Get(nil)
			Expect(request.Resource).To(Equal(def.ResourceT))
		})
	})
})
