package rest_test

import (
	. "github.com/gotgo/gokn/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SmallStruct struct {
	Message string
}

var _ = Describe("Rest", func() {

	It("should correctly create a local request", func() {
		cr := &ClientRequest{
			Body:     &SmallStruct{"hello"},
			Verb:     "POST",
			Resource: "/test",
		}
		req, resp := LocalRequest(cr)

		Expect(req.Raw.Method).To(Equal(cr.Verb))
		Expect(req.Raw.URL.Path).To(Equal(cr.Resource))
		target := req.Body.(*SmallStruct)
		Expect(target.Message).To(Equal(cr.Body.(*SmallStruct).Message))
		Expect(resp.Status).To(Equal(200))
	})

})
