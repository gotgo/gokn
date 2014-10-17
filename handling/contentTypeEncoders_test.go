package handling_test

import (
	. "github.com/krave-n/go/handling"
	"github.com/krave-n/go/rest"
	"github.com/gotgo/fw/tracing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type EncodeMe struct {
	Message string
}

var _ = Describe("ContentTypeEncoders", func() {

	var (
		encoders *ContentTypeEncoders
		response *rest.Response
		tracer   tracing.Tracer
	)

	BeforeEach(func() {
		tracer = &tracing.NopTracer{}
		encoders = NewContentTypeEncoders()
		response = &rest.Response{
			Body:        &EncodeMe{"test"},
			ContentType: "",
		}
	})

	It("should handle missing content type", func() {
		data := []byte{1, 2, 3}
		bts, err := encoders.Encode(data, "application/json")
		Expect(err).To(BeNil())
		Expect(bts).To(Equal(data))
	})

	It("should allow []byte to pass through without encoding, and without changing the content type", func() {
		data := []byte{1, 2, 3}
		bts, err := encoders.Encode(data, "application/json")
		Expect(err).To(BeNil())
		Expect(bts).To(Equal(data))
	})
})
