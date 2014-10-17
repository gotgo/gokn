package rest

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

type ResourceSpec struct {
	Definition           *ResourceDef
	requestContentTypes  []string //TODO: Request & Response ContentTypes go on the spec or the definition?
	responseContentTypes []string
	compiledTemplate     *template.Template
}

func NewResourceSpec(def *ResourceDef, reqContentTypes []string, respContentTypes []string) *ResourceSpec {
	return &ResourceSpec{
		Definition:           def,
		requestContentTypes:  reqContentTypes,
		responseContentTypes: respContentTypes,
	}
}

func (rsd *ResourceSpec) ResourceT() string {
	return rsd.Definition.ResourceT
}

func (rsd *ResourceSpec) ResourceArgs() interface{} {
	if rsd.Definition.ResourceArgs == nil {
		return nil
	} else {
		return reflect.New(rsd.Definition.ResourceArgs).Interface()
	}
}

func (rsd *ResourceSpec) Methods() []string {
	return rsd.Definition.Verbs
}

func (rsd *ResourceSpec) Headers() []string {
	return rsd.Definition.Headers
}

func (rsd *ResourceSpec) RequestBody() interface{} {
	if rsd.Definition.RequestBody == nil {
		return nil
	} else {
		//this returns an interface that is a pointer to a new instance of the datatype
		res := reflect.New(rsd.Definition.RequestBody).Interface()
		return res
	}
}

func (rsd *ResourceSpec) ResponseBody() interface{} {
	if rsd.Definition.ResponseBody == nil {
		return nil
	} else {
		//this returns an interface that is a pointer to a new instance of the datatype
		return reflect.New(rsd.Definition.ResponseBody).Interface()
	}
}

func (rsd *ResourceSpec) RequestContentTypes() []string {
	return rsd.requestContentTypes
}

func (rsd *ResourceSpec) ResponseContentTypes() []string {
	return rsd.responseContentTypes
}

// Client Behavior

func (rd *ResourceSpec) path(args interface{}) string {
	if args == nil {
		return rd.ResourceT()
	}

	rd.compile()
	toClean := structs.Map(args)
	clean := make(map[string]string)
	for k, v := range toClean {
		clean[k] = url.QueryEscape(fmt.Sprintf("%v", v))
	}
	buff := bytes.NewBufferString("")
	rd.compiledTemplate.Execute(buff, clean)
	return buff.String()
}

func (rd *ResourceSpec) compile() {
	if rd.compiledTemplate == nil {
		templateName := rd.ResourceT()
		resourceT := rd.resourceAsTemplate()
		compiled := template.Must(template.New(templateName).Parse(resourceT))
		rd.compiledTemplate = compiled
	}
}

// prepare converts the typical url template /url/{param} to the html.templates of
// /url/{{.param}}
func (rsd *ResourceSpec) resourceAsTemplate() string {
	urlT := rsd.ResourceT()
	halfway := strings.Replace(urlT, "{", "{{.", -1)
	final := strings.Replace(halfway, "}", "}}", -1)
	return final
}

func (rsd *ResourceSpec) Get(args interface{}) *ClientRequest {
	path := rsd.path(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "GET",
		Definition: rsd,
	}
	attachArgs(req, args)
	return req
}

func (rsd *ResourceSpec) Post(args interface{}, body interface{}) *ClientRequest {
	path := rsd.path(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "POST",
		Definition: rsd,
		Body:       body,
	}
	attachArgs(req, args)
	return req
}

func (rsd *ResourceSpec) Put(args interface{}, body interface{}) *ClientRequest {
	path := rsd.path(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "PUT",
		Definition: rsd,
		Body:       body,
	}
	attachArgs(req, args)
	return req
}

func (rsd *ResourceSpec) Patch(args interface{}, body interface{}) *ClientRequest {
	path := rsd.path(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "PATCH",
		Definition: rsd,
		Body:       body,
	}
	attachArgs(req, args)
	return req
}

func (rsd *ResourceSpec) Delete(args interface{}) *ClientRequest {
	path := rsd.path(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "Delete",
		Definition: rsd,
	}
	attachArgs(req, args)
	return req
}

func (rsd *ResourceSpec) Head(args interface{}) *ClientRequest {
	path := rsd.path(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "HEAD",
		Definition: rsd,
	}
	attachArgs(req, args)
	return req
}
