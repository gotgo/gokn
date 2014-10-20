package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/fatih/structs"
	"github.com/gotgo/fw/tracing"
)

// Client executes REST calls
type Client struct {
	Headers   map[string][]string
	Endpoints []ResourceEndpoint
	Encoder   func(v interface{}) ([]byte, error)
	Decoder   func(data []byte, v interface{}) error
	Tracer    tracing.RequestTracer
}

func NewClient() *Client {
	client := &Client{
		Encoder: json.Marshal,
		Decoder: json.Unmarshal,
		Tracer:  new(tracing.NopClientTracer),
	}
	return client
}

// SetHeader should receive a flat struct with no nested types. It will also overwrite existing headers.
func (t *Client) SetHeader(h interface{}) error {
	m := structs.Map(h)
	for k, v := range m {
		if structs.IsStruct(v) {
			return errors.New("Field of struct is a struct.  Can't set header value to a whole struct")
		}
		str := fmt.Sprintf("%v", v)
		t.Headers[k] = []string{str}
	}
	return nil
}

// Fetch makes it easier to Unmarshal the response from Do
func (c *Client) Fetch(r *ClientRequest, ctx *RequestContext, v interface{}) (*EndpointResponse, error) {
	if resp, err := c.Do(r, ctx); err != nil {
		return nil, err
	} else if bytes, err := resp.Bytes(); err != nil {
		return nil, err
	} else {
		err := c.unmarshal(bytes, &v)
		return resp, err
	}
}

func (c *Client) Do(r *ClientRequest, ctx *RequestContext) (*EndpointResponse, error) {
	tracer := ctx.Trace.NewRequest(resourceName(r), getArgs(r), r.Headers)
	tracer.Begin()
	defer tracer.End()

	client := new(http.Client)

	if req, err := c.newHttpRequest(r); err != nil {
		tracer.Annotate(tracing.Error, "request", err)
		return nil, err
	} else if resp, err := client.Do(req); err != nil {
		tracer.Annotate(tracing.Error, "request", err)
		return nil, err
	} else {
		resp := &EndpointResponse{
			HttpResponse: resp,
		}
		return resp, nil
	}
}

// newHttpRequest converts a ClientRequest into a http.Request
func (c *Client) newHttpRequest(cr *ClientRequest) (*http.Request, error) {
	resource, query := splitQueryPath(cr.Resource)

	var bts []byte
	if b, ok := cr.Body.([]byte); ok {
		bts = b
	} else if b, err := c.marshal(cr.Body); err != nil {
		return nil, err
	} else {
		bts = b
	}

	bodyCloser := ioutil.NopCloser(bytes.NewBuffer(bts))
	endpoint := c.endpoint() //TODO: retry on different endpoint if can't connect
	req := &http.Request{
		Method:        cr.Verb,
		Header:        cr.Headers,
		Body:          bodyCloser,
		ContentLength: int64(len(bts)),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		URL: &url.URL{
			Scheme:   endpoint.Scheme,
			Host:     endpoint.Host,
			Path:     path.Join(endpoint.ResourceRoot, resource),
			RawQuery: query,
		},
	}

	return req, nil
}

func splitQueryPath(r string) (path, query string) {
	parts := strings.Split(r, "?")
	path = parts[0]
	query = ""

	if len(parts) == 2 {
		query = parts[1]
	} else if len(parts) > 2 {
		panic("query portion of resource can't have more than one '?'")
	}
	return path, query
}

func resourceName(r *ClientRequest) string {
	if r.Definition == nil {
		return fmt.Sprintf("%s - %s", r.Verb, r.Resource)
	} else {
		return fmt.Sprintf("%s - %s", r.Verb, r.Definition.ResourceT)
	}
}

func (c *Client) endpoint() ResourceEndpoint {
	i := randInt(0, len(c.Endpoints)-1)
	return c.Endpoints[i]
}

func randInt(min int, max int) int {
	if min == 0 && max == 0 {
		return 0
	}
	return min + rand.Intn(max-min)
}

// marshal uses the Client.Encoder if it's not nil; otherwise,
// uses json.Marshal
func (t *Client) marshal(v interface{}) ([]byte, error) {
	e := t.Encoder
	if e == nil {
		e = json.Marshal
	}
	bytes, err := e(v)
	return bytes, err
}

// unmarshal uses the Client.Decoder if it's not nil; otherwise,
// uses json.Unmarsal
func (t *Client) unmarshal(bytes []byte, v interface{}) error {
	d := t.Decoder
	if d == nil {
		d = json.Unmarshal
	}
	err := d(bytes, &v)
	return err
}
