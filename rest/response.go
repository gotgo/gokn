package rest

import "reflect"

type Response struct {
	Body        interface{}
	Status      int
	Message     string
	Headers     map[string]string
	ContentType string
}

func (r *Response) IsStruct() bool {
	body := r.Body

	if body == nil {
		return false
	}

	t := reflect.TypeOf(body)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		return true
	} else {
		return false
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
