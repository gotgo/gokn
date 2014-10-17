package rest

// ClientRequest is ideal created from a ResourceSpec, although, there's nothing to prevent
// you from using it directly.
type ClientRequest struct {
	// Resource is the path to the resource not including the scheme, domain or port
	Resource string
	// One of the six Verbs aka Methods defined for RESTful calls
	Verb string
	// Definition is useful for metrics, monitoring & tracing
	Definition ServerResource
	// Body applies to POST, PUT & PATCH
	Body interface{}
	// Headers for this request.  If a header name already exists in the client,
	// these headers will override existing on the client;
	Headers map[string][]string
	// Args is the data used to combine with a ResourceT to produce the Resource
	args interface{}
}

func attachArgs(cr *ClientRequest, args interface{}) {
	cr.args = args
}

func getArgs(cr *ClientRequest) interface{} {
	return cr.args
}

// SetHeaders is a Fluent Method that Sets multiple headers on the existing ClientRequest headers,
// overwriting any header with the same name
func (cr *ClientRequest) SetHeaders(headers map[string][]string) *ClientRequest {
	if cr.Headers == nil {
		cr.Headers = make(map[string][]string)
	}
	for k, v := range headers {
		cr.Headers[k] = v
	}
	return cr
}

func (cr *ClientRequest) AddHeader(key, value string) *ClientRequest {
	if cr.Headers == nil {
		cr.Headers = make(map[string][]string)
	}

	list := cr.Headers[key]
	if list == nil {
		cr.Headers[key] = []string{value}
	} else {
		cr.Headers[key] = append(list, value)
	}
	return cr
}
