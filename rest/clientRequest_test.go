package rest_test

import (
	"github.com/gotgo/gokn/rest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientRequest", func() {

	var request *rest.ClientRequest

	BeforeEach(func() {
		request = &rest.ClientRequest{
			Resource: "/blah",
			Verb:     "GET",
		}
	})

	It("should set headers on a nil Request.Headers", func() {
		request.Headers = nil
		headers := make(map[string][]string)
		headers["a"] = []string{"yupA"}
		headers["b"] = []string{"yupB"}

		Expect(request.Headers).To(BeNil())

		request.SetHeaders(headers)
		Expect(len(request.Headers)).To(Equal(2))
		Expect(request.Headers["a"][0]).To(Equal(headers["a"][0]))
		Expect(request.Headers["b"][0]).To(Equal(headers["b"][0]))
	})

	It("should set headers on an existing Request.Headers and replace those with the same header name", func() {
		request.Headers = make(map[string][]string)
		request.Headers["x"] = []string{"yupX"}
		y := []string{"yupY"}
		request.Headers["y"] = y

		headers := make(map[string][]string)
		headers["a"] = []string{"yupA"}
		headers["b"] = []string{"yupB"}
		headers["x"] = []string{"replaced"}

		//starting expectation
		Expect(len(request.Headers)).To(Equal(2))

		//modify
		request.SetHeaders(headers)
		Expect(len(request.Headers)).To(Equal(4))
		Expect(request.Headers["a"][0]).To(Equal(headers["a"][0]))
		Expect(request.Headers["b"][0]).To(Equal(headers["b"][0]))
		Expect(request.Headers["x"][0]).To(Equal(headers["x"][0]))
		Expect(request.Headers["y"][0]).To(Equal(y[0]))
	})
})
