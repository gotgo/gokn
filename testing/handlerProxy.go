package testing

import "github.com/gotgo/gokn/rest"

func ConvertToHandler(args, body interface{}) (*rest.Request, *rest.Response) {
	req := &rest.Request{
		Raw:        nil,
		Context:    nil,
		Definition: nil,
		Args:       args,
		Body:       body,
	}

	//we start off in a success state
	resp := &rest.Response{
		Status: 200,
	}

	return req, resp
}
