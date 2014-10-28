package rest_test

import (
	"encoding/json"
	"reflect"

	. "github.com/gotgo/gokn/rest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServerResource", func() {
	var serverResource ServerResource

	BeforeEach(func() {
		def := &ResourceDef{
			ResourceT:    "/abc/{id}",
			ResourceArgs: nil,
			Verb:         "GET",
			Headers:      nil,
			RequestBody:  reflect.TypeOf(TestMessage{}),
		}

		ct := []string{"application/json"}
		serverResource = NewServerResource(def, ct, ct)
	})

	It("should be the correct instance type", func() {
		body := serverResource.RequestBody()

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
