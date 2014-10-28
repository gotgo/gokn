package rest

import "reflect"

// ResourceDef in a specification of the Resource.  Maybe rename to ResourceSpec
type ResourceDef struct {
	ResourceT    string // /sync/order
	ResourceArgs reflect.Type
	Verb         string   // GET POST
	Headers      []string //TODO: change to reflect.Type
	RequestBody  reflect.Type
	ResponseBody reflect.Type
	// where else would be put content type, if not here?
	RequestContentTypes  []string
	ResponseContentTypes []string
	template             *UrlPath
}

func (rd *ResourceDef) GetPath(args interface{}) string {
	template := rd.template
	if template == nil {
		rd.template = NewUrlPath(rd.ResourceT)
		template = rd.template
	}
	return template.Path(args)
}
