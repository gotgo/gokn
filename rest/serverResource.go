package rest

import "reflect"

// ServerResource is a Resource that the Server offers and is a readonly view
// into a ResourceDef
type ServerResource interface {
	//the Resource template
	ResourceT() string
	ResourceArgs() interface{}
	// Methods supported
	Verb() string
	// Headers Required
	Headers() []string
	// Request returns a new instance of the request
	RequestBody() interface{}
	// Response return a new instance of the response
	ResponseBody() interface{}
	// RequestContentTypes are a list of accepted content-type for the request portion
	RequestContentTypes() []string
	// ResponseContentTypes are a list content-type that the response will be sent in
	ResponseContentTypes() []string //not sure this should be an array?
}

func NewServerResource(definition *ResourceDef, reqContentTypes []string, respContentTypes []string) ServerResource {
	return &serverResourceSpec{
		Definition:           definition,
		requestContentTypes:  reqContentTypes,
		responseContentTypes: respContentTypes,
	}
}

type serverResourceSpec struct {
	Definition           *ResourceDef
	requestContentTypes  []string //TODO: Request & Response ContentTypes go on the spec or the definition?
	responseContentTypes []string
}

func (rsd *serverResourceSpec) ResourceT() string {
	return rsd.Definition.ResourceT
}

func (rsd *serverResourceSpec) ResourceArgs() interface{} {
	if rsd.Definition.ResourceArgs == nil {
		return nil
	} else {
		return reflect.New(rsd.Definition.ResourceArgs).Interface()
	}
}

func (rsd *serverResourceSpec) Verb() string {
	return rsd.Definition.Verb
}

func (rsd *serverResourceSpec) Headers() []string {
	return rsd.Definition.Headers
}

func (rsd *serverResourceSpec) RequestBody() interface{} {
	if rsd.Definition.RequestBody == nil {
		return nil
	} else {
		//this returns an interface that is a pointer to a new instance of the datatype
		res := reflect.New(rsd.Definition.RequestBody).Interface()
		return res
	}
}

func (rsd *serverResourceSpec) ResponseBody() interface{} {
	if rsd.Definition.ResponseBody == nil {
		return nil
	} else {
		//this returns an interface that is a pointer to a new instance of the datatype
		return reflect.New(rsd.Definition.ResponseBody).Interface()
	}
}

func (rsd *serverResourceSpec) RequestContentTypes() []string {
	return rsd.requestContentTypes
}

func (rsd *serverResourceSpec) ResponseContentTypes() []string {
	return rsd.responseContentTypes
}
