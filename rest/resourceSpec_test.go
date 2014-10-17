package rest_test

import (
	"encoding/json"
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
			Verbs:        []string{"GET"},
			Headers:      nil,
			RequestBody:  reflect.TypeOf(TestMessage{}),
		}

		spec = &rest.ResourceSpec{Definition: def}
	})

	Context("server behavior", func() {

		It("should be the correct instance type", func() {
			body := spec.RequestBody()

			bodyTyped, ok := body.(*TestMessage)

			Expect(ok).To(Equal(true))

			bodyTyped.Message = "haha"

			tm := &TestMessage{Message: "test message"}
			bts, err := json.Marshal(tm)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bts, &body)
			if err != nil {
				panic(err)
			}

			Expect(ok).To(Equal(true))
			Expect(bodyTyped.Message).To(Equal(tm.Message))
		})
	})

	Context("client behavior", func() {
		It("should handle a nil argument, and return the spec with it's template", func() {
			var cl rest.ClientResource
			cl = spec

			request := cl.Get(nil)
			Expect(request.Resource).To(Equal(spec.ResourceT()))
		})
	})
})
