package rest_test

import (
	"encoding/json"
	"fmt"

	"github.com/gotgo/fw/io"
	"github.com/gotgo/gokn/rest"
	"github.com/gotgo/gokn/testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Client Test relies on an echo service runing on localhost:23000
// There is an echo service the the testing folder that can be run to satisfy this
// requirement
var (
	listener *io.GracefulListener
)

var _ = BeforeSuite(func() {
	fmt.Println("******* this suite requires the testing echo service to be running on 23000 *****")

	listener = testing.StartEchoHandler()
})

var _ = AfterSuite(func() {
	listener.Shutdown()
})

type echo struct {
	Message string `json:"message"`
}

var _ = Describe("Client", func() {

	var (
		client  *rest.Client
		context *rest.RequestContext
	)

	BeforeEach(func() {
		endpoint := &rest.ResourceEndpoint{
			Host:   "localhost:23000",
			Scheme: "http",
		}
		client = &rest.Client{
			Endpoints: []rest.ResourceEndpoint{*endpoint},
		}
		context = rest.NewRequestContext()
	})

	Context("Client DO", func() {
		AssertDoExpectations := func(req *rest.ClientRequest) {
			resp, err := client.Do(req, context)

			hr := resp.HttpResponse

			Expect(hr.StatusCode).To(Equal(200))
			Expect(err).To(BeNil())

			Expect(resp).ToNot(BeNil())
			bytes, err := resp.Bytes()
			Expect(err).To(BeNil())

			Expect(hr.ContentLength).To(BeNumerically(">", 200))
			er := &testing.EchoRequest{}
			err = json.Unmarshal(bytes, &er)
			Expect(err).To(BeNil())
			Expect(len(bytes)).To(BeNumerically(">", 1))
			if req.Verb == "POST" || req.Verb == "PUT" {
				mp := er.Body.(map[string]interface{})
				Expect(mp["message"]).To(Equal("test message"))
			}
		}
		It("Should GET", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "GET",
			}
			AssertDoExpectations(req)
		})

		It("Should POST", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "POST",
				Body:     &echo{"test message"},
			}
			AssertDoExpectations(req)
		})
		It("Should PUT", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "PUT",
				Body:     &echo{"test message"},
			}
			AssertDoExpectations(req)
		})
		It("Should DELETE", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "DELETE",
			}
			AssertDoExpectations(req)
		})
		It("Should HEAD ", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "HEAD",
			}
			resp, err := client.Do(req, context)

			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(err).To(BeNil())
			//js, err := json.MarshalIndent(resp, "", "\t")
			//os.Stdout.Write(js)

		})
		It("Should PATCH", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "PATCH",
			}
			AssertDoExpectations(req)
		})
	})

	Context("Client Fetch", func() {
		AssertDoExpectations := func(req *rest.ClientRequest) {
			echoReq := &testing.EchoRequest{}
			resp, err := client.Fetch(req, context, echoReq)

			hr := resp.HttpResponse

			Expect(hr.StatusCode).To(Equal(200))
			Expect(err).To(BeNil())

			_, err = resp.Bytes()
			Expect(err).ToNot(BeNil())

			Expect(hr.ContentLength).To(BeNumerically(">", 200))

			if req.Verb == "POST" || req.Verb == "PUT" {
				mp := echoReq.Body.(map[string]interface{})
				Expect(mp["message"]).To(Equal("test message"))
			}
		}
		It("Should GET", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "GET",
			}
			AssertDoExpectations(req)
		})

		It("Should POST", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "POST",
				Body:     &echo{"test message"},
			}
			AssertDoExpectations(req)
		})
		It("Should PUT", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "PUT",
				Body:     &echo{"test message"},
			}
			AssertDoExpectations(req)
		})
		It("Should DELETE", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "DELETE",
			}
			AssertDoExpectations(req)
		})
		It("Should HEAD ", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "HEAD",
			}
			resp, err := client.Do(req, context)

			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(err).To(BeNil())
			//js, err := json.MarshalIndent(resp, "", "\t")
			//os.Stdout.Write(js)

		})
		It("Should PATCH", func() {
			req := &rest.ClientRequest{
				Resource: "",
				Verb:     "PATCH",
			}
			AssertDoExpectations(req)
		})
	})

})
