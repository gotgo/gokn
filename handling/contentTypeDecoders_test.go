package handling_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	. "github.com/gotgo/gokn/handling"
	"github.com/gotgo/gokn/rest"
	"github.com/gotgo/fw/tracing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ByteAlias []byte

type RequestBody struct {
	Message string `json:"message"`
}

var _ = Describe("ContentTypeDecoder", func() {

	var (
		decoders    *ContentTypeDecoders
		request     *rest.Request
		requestBody *RequestBody
		resource    *rest.ResourceDef
		tracer      tracing.Tracer
	)

	var setBody = func(bts []byte, err error) {
		if err != nil {
			panic(err)
		}
		request.Raw.Body = ioutil.NopCloser(bytes.NewReader(bts))
	}

	BeforeEach(func() {
		decoders = NewContentTypeDecoders()
		tracer = &tracing.NopTracer{}
		resource = &rest.ResourceDef{
			ResourceT:   "/decodeMe",
			RequestBody: reflect.TypeOf(RequestBody{}),
			Headers:     []string{"GET"},
		}

		request = &rest.Request{
			Definition: &rest.ResourceSpec{Definition: resource},
			Raw: &http.Request{
				Header: make(http.Header),
			},
		}

		requestBody = &RequestBody{Message: "test message"}
	})

	It("should default to json", func() {

		request.Raw.Header["Content-Type"] = []string{""}
		setBody(json.Marshal(requestBody))
		err := decoders.DecodeBody(request, tracer)
		Expect(err).To(BeNil())

		v, ok := request.Body.(*RequestBody)
		Expect(ok).To(Equal(true))
		Expect(v.Message).To(Equal(requestBody.Message))
	})
	It("should decode json", func() {

		request.Raw.Header["Content-Type"] = []string{"application/json"}
		setBody(json.Marshal(requestBody))
		err := decoders.DecodeBody(request, tracer)
		Expect(err).To(BeNil())

		v, ok := request.Body.(*RequestBody)
		Expect(ok).To(Equal(true))
		Expect(v.Message).To(Equal(requestBody.Message))
	})
	It("should not decode if the expected type is []byte", func() {
		resource.RequestBody = reflect.TypeOf(ByteAlias{})
		request.Raw.Header["Content-Type"] = []string{"application/json"}
		data := []byte{1, 3, 5, 4}
		setBody(data, nil)
		decoders.DecodeBody(request, tracer)
		Expect(request.Body).To(Equal(data))
	})

})
