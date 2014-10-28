package testing

import "github.com/gotgo/gokn/rest"

func ConvertToHandler(args, body interface{}, resource ServerResource) (*rest.Request, rest.Responder) {
	req := &rest.Request{
		Raw:        nil,
		Context:    nil,
		Definition: resource,
		Args:       args,
		Body:       body,
	}

	resp := &rest.Response{}

	return req, resp
}
