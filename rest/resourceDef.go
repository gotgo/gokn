package rest

import (
	"reflect"
	"text/template"
)

// ResourceDef in a specification of the Resource.  Maybe rename to ResourceSpec
type ResourceDef struct {
	ResourceT    string // /sync/order
	ResourceArgs reflect.Type
	Verbs        []string // GET POST
	Headers      []string
	RequestBody  reflect.Type
	ResponseBody reflect.Type
	// where else would be put content type, if not here?
	RequestContentTypes  []string
	ResponseContentTypes []string
	compiledTemplate     *template.Template
}
