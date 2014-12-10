package rest

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/oleiade/reflections"
)

//more work to support http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

func GetHeaderValue(key string, headers map[string][]string) string {
	values := headers[key]
	if values != nil && len(values) > 0 {
		return values[0]
	}
	return ""
}

// currently used primarily for testing, however, could be used to make request to the
// local process for cases where it made more sense to delpoy a typically remote resource locally
func LocalRequest(cr *ClientRequest) (*Request, *Response) {
	response := NewResponse()
	client := NewClient()
	if rawReq, err := client.NewHttpRequest(cr); err != nil {
		panic(err)
	} else {
		request := NewRequest(rawReq, NewRequestContext(), cr.Definition)
		request.Body = cr.Body
		request.Args = cr.args
		return request, response
	}
}

func FullApiWithHandlers(apiSpec interface{}) map[ServerResource]Handler {
	handlers := make(map[ServerResource]Handler)

	items, err := reflections.Items(apiSpec)
	if err != nil {
		panic(err)
	}

	for _, v := range items {
		if target, ok := v.(*ResourceSpec); ok {
			verbs, handler := target.ServeAll()
			if handler == nil {
				panic(fmt.Sprintf("no handler for endpoint %v ", target))
			}
			for _, verb := range verbs {
				handlers[verb] = handler
			}
		}
	}
	return handlers
}

// Bytes, is a helper method to reduce the number of lines to get a byte array out of the
// EndpointResponse
func Bytes(resp *EndpointResponse, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("no response from endpoint or code incorrectly called.")
	}

	if resp.HttpResponse == nil {
		return nil, errors.New("No HttpResponse Response")
	}
	if resp.HttpResponse.StatusCode != http.StatusOK {
		return nil, errors.New(resp.HttpResponse.Status)
	}

	if bytes, err := resp.Bytes(); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

// Encode - Returns a URL that is encoded in a way that makes sense, preseves as many special characters
// as possible. The default encoding from GO makes no sense as it encodes even valid values such as + , :
func Encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	path := buf.String()
	u := &url.URL{
		Path: path,
	}
	return u.RequestURI()
}
