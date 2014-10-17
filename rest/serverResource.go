package rest

// ServerResource is a Resource that the Server offers and is a readonly view
// into a ResourceDef
type ServerResource interface {
	//the Resource template
	ResourceT() string
	ResourceArgs() interface{}
	// Methods supported
	Methods() []string
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
