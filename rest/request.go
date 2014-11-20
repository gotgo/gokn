package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/gotgo/fw/tracing"
	"github.com/gotgo/fw/util"
)

type Request struct {
	Raw        *http.Request
	Context    *RequestContext
	Definition ServerResource
	Args       interface{}
	Body       interface{}
	//header????
	bodyBytes []byte
}

func NewRequest(raw *http.Request, ctx *RequestContext, spec ServerResource) *Request {
	r := &Request{
		Raw:        raw,
		Context:    ctx,
		Definition: spec,
	}
	return r
}

func (r *Request) ContentType() string {
	if r.Raw != nil {
		ct := r.Raw.Header["Content-Type"]
		if len(ct) > 0 {
			return ct[0]
		}
	}
	return ""
}

func (r *Request) Annotate(f tracing.From, k string, v interface{}) {
	r.Context.Trace.Annotate(f, k, v)
}

// Bytes returns the body of the request as a []byte
func (r *Request) Bytes() ([]byte, error) {
	if r.bodyBytes == nil {
		defer r.Raw.Body.Close()

		if bts, err := ioutil.ReadAll(r.Raw.Body); err != nil {
			return nil, err
		} else {
			r.bodyBytes = bts
		}
	}
	return r.bodyBytes, nil
}

func (r *Request) DecodeArgs(argValues map[string]string) error {
	if args := r.Definition.ResourceArgs(); args != nil && len(argValues) > 0 {
		if err := util.MapToStruct(argValues, &args); err != nil {
			return err
		} else {
			r.Args = args
		}
	}
	return nil
}
