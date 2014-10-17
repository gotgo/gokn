package handling_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	. "github.com/gotgo/gokn/handling"
	"github.com/gotgo/gokn/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestResponseWriter struct {
	WriteReturnCount int
	WriteReturnError error
	WriteBytes       []byte
	WriteHeaderCode  int
}

func (trw *TestResponseWriter) Header() http.Header {
	return make(map[string][]string)
}

func (trw *TestResponseWriter) Write(bytes []byte) (int, error) {
	trw.WriteBytes = bytes
	return trw.WriteReturnCount, trw.WriteReturnError
}

func (trw *TestResponseWriter) WriteHeader(code int) {
	trw.WriteHeaderCode = code
}

type TestRouter struct {
	RegisterCount int
	GetCount      int
	PostCount     int
	PutCount      int
	DeleteCount   int
	HeadCount     int
	PatchCount    int
	Handlers      []func(http.ResponseWriter, *http.Request)
}

func NewTestRouter() *TestRouter {
	router := new(TestRouter)
	router.Handlers = []func(http.ResponseWriter, *http.Request){}
	return router
}

func (tr *TestRouter) RequestArgs(req *http.Request) map[string]string {
	return make(map[string]string)
}

func (tr *TestRouter) RegisterRoute(verb, path string, f func(http.ResponseWriter, *http.Request)) {
	tr.RegisterCount++
	switch verb {
	case "GET":
		tr.GetCount++
	case "POST":
		tr.PostCount++
	case "PUT":
		tr.PutCount++
	case "DELETE":
		tr.DeleteCount++
	case "HEAD":
		tr.HeadCount++
	case "PATCH":
		tr.PatchCount++
	}
	tr.Handlers = append(tr.Handlers, f)
}

type TestHandler struct {
	ResponseStatus int
}
type TestStruct struct {
	Message string
}

func NewTestHandler() *TestHandler {
	h := new(TestHandler)
	h.ResponseStatus = 200
	return h
}

func (th *TestHandler) setResponse(resp rest.Responder) {
	if th.ResponseStatus != 200 {
		resp.SetStatus(th.ResponseStatus, "")
	}

}

func (th *TestHandler) Get(req *rest.Request, resp rest.Responder) {
	th.setResponse(resp)
}
func (th *TestHandler) Post(req *rest.Request, resp rest.Responder) {
	resp.SetBody(&TestStruct{"response"})
	th.setResponse(resp)
}
func (th *TestHandler) Put(req *rest.Request, resp rest.Responder) {
	resp.SetBody(&TestStruct{"response"})
	th.setResponse(resp)
}
func (th *TestHandler) Delete(req *rest.Request, resp rest.Responder) {
	th.setResponse(resp)
}
func (th *TestHandler) Head(req *rest.Request, resp rest.Responder) {
	th.setResponse(resp)
}
func (th *TestHandler) Patch(req *rest.Request, resp rest.Responder) {
	th.setResponse(resp)
}

func getSpec(resourceT string, verbs ...string) rest.ServerResource {
	def := &rest.ResourceDef{
		ResourceT:    resourceT,
		Verbs:        verbs,
		RequestBody:  reflect.TypeOf([]byte{}),
		ResponseBody: reflect.TypeOf([]byte{}),
	}
	return &rest.ResourceSpec{Definition: def}
}

var _ = Describe("RootHandler", func() {

	var (
		root    *RootHandler
		handler *TestHandler
		router  *TestRouter
		request *http.Request
		writer  *TestResponseWriter
	)

	BeforeEach(func() {
		root = NewRootHandler()
		handler = NewTestHandler()
		router = NewTestRouter()

		tm := &TestStruct{Message: ""}

		buf, _ := json.Marshal(tm)
		request = &http.Request{
			Method:        "POST",
			Body:          ioutil.NopCloser(bytes.NewReader(buf)),
			ContentLength: int64(len(buf)),
		}

		writer = new(TestResponseWriter)
	})

	Context("Test Init", func() {
		It("should have reasonable defaults set after NewRootHandler", func() {
			Expect(root.Log).ToNot(BeNil())
			Expect(root.Binder).ToNot(BeNil())
		})
	})

	Context("Bind", func() {
		It("should bind GET", func() {
			spec := getSpec("/test", "GET")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(1))
			Expect(router.GetCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(1))
		})
		It("should bind POST", func() {
			spec := getSpec("/test", "POST")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(1))
			Expect(router.PostCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(1))
		})
		It("should bind PUT", func() {
			spec := getSpec("/test", "PUT")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(1))
			Expect(router.PutCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(1))
		})
		It("should bind DELETE", func() {
			spec := getSpec("/test", "DELETE")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(1))
			Expect(router.DeleteCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(1))
		})
		It("should bind HEAD", func() {
			spec := getSpec("/test", "HEAD")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(1))
			Expect(router.HeadCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(1))
		})
		It("should bind PATCH", func() {
			spec := getSpec("/test", "PATCH")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(1))
			Expect(router.PatchCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(1))
		})
		It("should bind to all verbs in one go", func() {
			spec := getSpec("/test", "GET", "POST", "PUT", "DELETE", "HEAD", "PATCH")
			root.Bind(router, spec, handler, "")
			Expect(router.RegisterCount).To(Equal(6))
			Expect(router.GetCount).To(Equal(1))
			Expect(router.PostCount).To(Equal(1))
			Expect(router.PutCount).To(Equal(1))
			Expect(router.DeleteCount).To(Equal(1))
			Expect(router.HeadCount).To(Equal(1))
			Expect(router.PatchCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(6))
		})
	})

	Context("BindAll", func() {
		It("should work", func() {
			spec1 := getSpec("/testA", "GET")
			spec2 := getSpec("/testB", "GET")
			spec3 := getSpec("/testC", "POST")
			specs := make(map[rest.ServerResource]rest.Handler)
			specs[spec1] = handler
			specs[spec2] = handler
			specs[spec3] = handler
			root.BindAll(router, specs, "")
			Expect(router.RegisterCount).To(Equal(3))
			Expect(router.GetCount).To(Equal(2))
			Expect(router.PostCount).To(Equal(1))
			Expect(len(router.Handlers)).To(Equal(3))
		})
	})

	Context("Calling wrapped handler", func() {

		FIt("should write bytes", func() {
			spec := getSpec("/test", "POST")
			root.Bind(router, spec, handler, "")
			wrappedHandler := router.Handlers[0]

			request.Header = make(map[string][]string)
			request.Header["Content-Type"] = []string{"application/json"}
			wrappedHandler(writer, request)
			Expect(writer.WriteBytes).ToNot(BeNil())
			Expect(len(writer.WriteBytes)).Should(BeNumerically(">", 0))
		})

		It("should write header value of error on response error", func() {
			spec := getSpec("/test", "POST")
			root.Bind(router, spec, handler, "")
			wrappedHandler := router.Handlers[0]

			status := 401
			handler.ResponseStatus = status
			wrappedHandler(writer, request)
			Expect(len(writer.WriteBytes)).To(Equal(0))
			Expect(writer.WriteHeaderCode).To(Equal(status))
		})

		It("should return an error code on Encode failure", func() {

			root.Encoders.Set(&ContentTypeEncoder{
				ContentType: "fail",
				Encode: func(v interface{}) ([]byte, error) {
					return nil, errors.New("fail")
				},
			})
			def := &rest.ResourceDef{
				ResourceT: "willfail",
				Verbs:     []string{"GET"},
			}

			contentType := []string{"fail"}
			spec := rest.NewResourceSpec(def, contentType, contentType)
			root.Bind(router, spec, handler, "")
			router.Handlers[0](writer, request)
			Expect(writer.WriteHeaderCode).To(Equal(http.StatusInternalServerError))
		})

	})

})
