package handling

import (
	"bytes"
	"errors"
	"fmt"

	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gotgo/fw/logging"
	"github.com/gotgo/fw/me"
	"github.com/gotgo/fw/tracing"
	"github.com/gotgo/gokn/rest"
)

// RootHandler binds api endpoint to a router
//
//	Example:
//
//		func newSetupHandlers(router *mux.Router, graph inject.Graph) {
//			//Set Custom Binder
//			kandle := new(KraveHandler)
//			root := &RootHandler{
//				Binder: kandle.RequiresDeviceAuthentication,
//			}
//
//			//Bind
//			pingEndpoint := new(PingEndpoint)
//			pingHandler := new(PingHandler)
//			root.Bind(router, pingEndpoint, pingHandler)
//		}
//
type RootHandler struct {
	Log          logging.Logger `inject:""`
	Binder       BindingFunc
	TraceHeader  string
	SpanHeader   string
	Encoders     *ContentTypeEncoders
	Decoders     *ContentTypeDecoders
	TraceHandler func(*tracing.TraceMessage)
}

func NewRootHandler() *RootHandler {
	root := &RootHandler{
		Log:          new(logging.NoOpLogger),
		Binder:       AnonymousHandler,
		TraceHeader:  traceHeader,
		SpanHeader:   spanHeader,
		Encoders:     NewContentTypeEncoders(),
		Decoders:     NewContentTypeDecoders(),
		TraceHandler: func(*tracing.TraceMessage) {},
	}

	return root
}

const (
	traceHeader = "tr-trace"
	spanHeader  = "tr-span"
)

func (rh *RootHandler) convertRequestResponse(w http.ResponseWriter, r *http.Request, endpoint rest.ServerResource) (*rest.Request, *rest.Response) {

	request := rest.NewRequest(r, rest.NewRequestContext(), endpoint)

	response := &rest.Response{
		Status:  200,
		Message: "ok",
		Headers: make(map[string]string),
	}
	return request, response
}

func requestName(r *rest.Request) string {
	rs := r.Definition.ResourceT()
	return fmt.Sprintf("%s - %s", r.Raw.Method, rs)
}

func setResponseContentType(response *rest.Response, req *http.Request, resp http.ResponseWriter, endpoint rest.ServerResource) {
	// TODO: ideally we'd match the preferred accept type, with a type we can respond with,
	// for now, until we support this, just ignore the accept string and return the first
	// content type we know we can return
	if response.ContentType == "" {
		cts := endpoint.ResponseContentTypes()
		if cts != nil && len(cts) > 0 {
			response.ContentType = cts[0]
		} else if cts = req.Header["Content-Type"]; cts != nil && len(cts) > 0 {
			response.ContentType = cts[0] //try returning the same as the requested type
		} else if cts = endpoint.RequestContentTypes(); cts != nil && len(cts) > 0 {
			response.ContentType = cts[0]
		}
	}
	resp.Header()["Content-Type"] = []string{response.ContentType}
}

type responseData struct {
	Data          []byte
	StatusCode    int
	StatusMessage string
	PanicMessage  string
}

func getErrorMessage(e interface{}) string {
	var msg string
	if err, ok := e.(error); ok {
		msg = err.Error()
	} else if str, ok := e.(string); ok {
		msg = str
	} else if _, ok := e.(runtime.Error); ok {
		msg = err.Error()
	} else {
		msg = ""
	}
	return msg
}

func (root *RootHandler) guaranteedReply(writer http.ResponseWriter, response *responseData, trace *tracing.TraceMessage) {
	defer root.TraceHandler(trace)

	var panicMessage string
	if r := recover(); r != nil {
		stack := make([]byte, 2048)
		runtime.Stack(stack, true)
		panicMessage = getErrorMessage(r)
		response.StatusMessage = "Internal Server Error"
		response.StatusCode = 500
		stackTrace := fmt.Sprintf("%s callstack: %s", panicMessage, stack)
		trace.Annotate(tracing.FromPanic, "request fail", panicMessage)
		trace.Annotate(tracing.FromPanic, "stack", stackTrace)
		root.Log.Error("Panic Occured", me.NewErr(stackTrace))
	}

	if response.StatusCode != 200 {
		if response.StatusCode == 0 {
			response.StatusCode = 500
			if response.StatusMessage == "" {
				response.StatusMessage = "Internal Server Error: Failed to complete"
			}
			root.Log.Error("Unhandled Panic", errors.New("Request Not completed"))
		}

		trace.Annotate(tracing.FromError, fmt.Sprintf("httpResponse: %v", response.StatusCode), response.StatusMessage)
		trace.RequestFail()
		http.Error(writer, response.StatusMessage, response.StatusCode)
		writer.Write([]byte{})
	} else {
		data := response.Data
		if data == nil {
			data = []byte{}
		}

		if bytesSent, err := writer.Write(data); err != nil {
			root.Log.Warn("failed to write response",
				&logging.KV{"message", "partial reply, failed to send entire reply"},
				&logging.KV{"bytesSent", bytesSent},
				&logging.KV{"totalBytes", len(data)},
			)
		}
		trace.RequestCompleted()
	}
}

func flattenForm(form map[string][]string) map[string]string {
	m := make(map[string]string)
	for k, v := range form {
		m[k] = strings.Join(v, ",")
	}
	return m
}

// TODO: remove direct dependency on gorilla mux, for path
func (root *RootHandler) hackReqArgs(req *http.Request, existingArgs map[string]string) {
	//get the /avar/{avar}/bvar/{bvar}
	pa := mux.Vars(req)
	if pa != nil {
		//TODO: right now this overwrites... maybe it should not??
		for k, v := range pa {
			existingArgs[k] = v
		}
	}
}

func (root *RootHandler) createHttpHandler(handler rest.HandlerFunc, endpoint rest.ServerResource) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		traceUid := rest.GetHeaderValue(root.TraceHeader, r.Header)
		spanUid := rest.GetHeaderValue(root.SpanHeader, r.Header)
		traceMessage := tracing.NewReceiveTrace(traceUid, spanUid)
		tracer := tracing.NewMessageTracer(traceMessage)
		responseData := &responseData{}
		defer root.guaranteedReply(w, responseData, traceMessage)

		//should ParseMultipartForm be configurable?? so it's only called when needed?
		r.ParseMultipartForm(120000)
		args := flattenForm(r.Form)
		root.hackReqArgs(r, args)

		request, response := root.convertRequestResponse(w, r, endpoint)
		request.Context.Trace = tracer

		traceMessage.ReceivedRequest(requestName(request), args, r.Header)

		if err := request.DecodeArgs(args); err != nil {
			responseData.StatusCode = http.StatusBadRequest
			responseData.StatusMessage = "Bad Request: failed parse expected URL parameters"
			return
		}

		if err := root.Decoders.DecodeBody(request, traceMessage); err != nil {
			responseData.StatusCode = http.StatusBadRequest
			responseData.StatusMessage = "Bad Request: Failed to decode request body for the provided Content-Type"
			return
		}

		boundHandler := root.Binder(handler)
		boundHandler(request, response)

		if response.Error != nil {
			request.Context.Trace.Annotate(tracing.FromError, "request failed", response.Error)
		}

		responseData.StatusCode = response.Status

		if response.Status != http.StatusOK {
			responseData.StatusMessage = response.Message
			return
		}

		setResponseContentType(response, r, w, endpoint)

		var bts []byte
		var err error

		if bts, err = root.Encoders.Encode(response.Body, response.ContentType); err != nil {
			responseData.StatusCode = http.StatusInternalServerError
			responseData.StatusMessage = "Internal Server Error - Failed to encode response body"
			return
		}
		w.Header()["Content-Length"] = []string{strconv.Itoa(len(bts))}
		responseData.Data = bts

		if response.IsBinary() {
			traceMessage.AnnotateBinary(tracing.FromResponseData, "body", bytes.NewReader(bts), response.ContentType)
		} else {
			traceMessage.Annotate(tracing.FromResponseData, "body", fmt.Sprintf("%s", response.Body))
		}
	}
}

func (root *RootHandler) Bind(router SimpleRouter, endpoint rest.ServerResource, handler rest.Handler, resourceRoot string) {
	if handler == nil {
		panic(fmt.Sprintf("handler can't be nil", endpoint))
	}

	resourcePathT := path.Join(resourceRoot, endpoint.ResourceT())
	httpMethod := endpoint.Verb()
	errMessage := "can't bind. method named %s is missing from type %s"
	handlerName := reflect.TypeOf(handler).Name()
	var fn rest.HandlerFunc

	if httpMethod == "GET" {
		if h, ok := handler.(rest.GetHandler); !ok {
			panic(fmt.Sprintf(errMessage, httpMethod, handlerName, resourcePathT))
		} else {
			fn = h.Get
		}
	} else if httpMethod == "POST" {
		if h, ok := handler.(rest.PostHandler); !ok {
			panic(fmt.Sprintf(errMessage, httpMethod, handlerName, resourcePathT))
		} else {
			fn = h.Post
		}
	} else if httpMethod == "PUT" {
		if h, ok := handler.(rest.PutHandler); !ok {
			panic(fmt.Sprintf(errMessage, httpMethod, handlerName, resourcePathT))
		} else {
			fn = h.Put
		}
	} else if httpMethod == "DELETE" {
		if h, ok := handler.(rest.DeleteHandler); !ok {
			panic(fmt.Sprintf(errMessage, httpMethod, handlerName, resourcePathT))
		} else {
			fn = h.Delete
		}
	} else if httpMethod == "HEAD" {
		if h, ok := handler.(rest.HeadHandler); !ok {
			panic(fmt.Sprintf(errMessage, httpMethod, handlerName, resourcePathT))
		} else {
			fn = h.Head
		}
	} else if httpMethod == "PATCH" {
		if h, ok := handler.(rest.PatchHandler); !ok {
			panic(fmt.Sprintf(errMessage, httpMethod, handlerName, resourcePathT))
		} else {
			fn = h.Patch
		}
	}

	wrappedHandler := root.createHttpHandler(fn, endpoint)
	router.RegisterRoute(httpMethod, resourcePathT, wrappedHandler)
	root.Log.Inform(fmt.Sprintf("Bound endpoint %s %s", httpMethod, resourcePathT))
}

// BindAll is a helper for calling Bind on a list of endpoints
func (root *RootHandler) BindAll(router SimpleRouter, endpoints map[rest.ServerResource]rest.Handler, resourceRoot string) {
	for definition, handler := range endpoints {
		root.Bind(router, definition, handler, resourceRoot)
	}
}
