package rest

import "reflect"

type Response struct {
	Body        interface{}
	Status      int
	Message     string
	Headers     map[string]string
	ContentType string
}

func NewResponse() *Response {
	return &Response{
		Status:  200,
		Message: "OK",
		Headers: make(map[string]string),
	}
}

func (r *Response) IsBinary() bool {
	body := r.Body

	if body == nil {
		return false
	}

	if _, ok := body.(string); ok {
		return false
	}

	if _, ok := body.([]byte); ok {
		return true
	}

	t := reflect.TypeOf(body)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		return false
	} else {
		return true
	}
}

func (r *Response) SetBody(data interface{}) {
	r.Body = data
}

func (r *Response) SetStatus(status int, message string) {
	r.Status = status
	r.Message = message
}

func (r *Response) SetContentType(contentType string) {
	r.ContentType = contentType
}

func (r *Response) AddHeader(key, value string) {
	r.Headers[key] = value
}
