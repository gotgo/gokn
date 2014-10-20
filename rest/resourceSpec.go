package rest

type ResourceSpec struct {
	defaultContentType []string
	get                *ResourceDef
	put                *ResourceDef
	post               *ResourceDef
	delete             *ResourceDef
	patch              *ResourceDef
	head               *ResourceDef
}

func NewResourceSpec(defaultContentType string) *ResourceSpec {
	return &ResourceSpec{
		defaultContentType: []string{defaultContentType},
	}
}

func (r *ResourceSpec) Use(def *ResourceDef) *ResourceSpec {
	switch def.Verb {
	case "GET":
		r.get = def
	case "POST":
		r.post = def
	case "PUT":
		r.put = def
	case "DELETE":
		r.delete = def
	case "HEAD":
		r.head = def
	case "PATCH":
		r.patch = def
	}
	return r
}

func (rs *ResourceSpec) ServeAll() []ServerResource {
	ct := rs.defaultContentType
	all := make([]ServerResource, 0)
	if rs.get != nil {
		all = append(all, NewServerResource(rs.get, ct, ct))
	}

	if rs.post != nil {
		all = append(all, NewServerResource(rs.post, ct, ct))
	}

	if rs.put != nil {
		all = append(all, NewServerResource(rs.put, ct, ct))
	}

	if rs.delete != nil {
		all = append(all, NewServerResource(rs.delete, ct, ct))
	}

	if rs.head != nil {
		all = append(all, NewServerResource(rs.head, ct, ct))
	}

	if rs.patch != nil {
		all = append(all, NewServerResource(rs.patch, ct, ct))
	}
	return all
}

// Client Behavior

func (rs *ResourceSpec) Get(args interface{}) *ClientRequest {
	if rs.get == nil {
		panic("get is nil")
	}

	path := rs.get.GetPath(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "GET",
		Definition: NewServerResource(rs.get, rs.defaultContentType, rs.defaultContentType),
	}
	attachArgs(req, args)
	return req
}

func (rs *ResourceSpec) Post(args interface{}, body interface{}) *ClientRequest {
	if rs.post == nil {
		panic("post is nil")
	}

	path := rs.post.GetPath(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "POST",
		Definition: NewServerResource(rs.post, rs.defaultContentType, rs.defaultContentType),
		Body:       body,
	}
	attachArgs(req, args)
	return req
}

func (rs *ResourceSpec) Put(args interface{}, body interface{}) *ClientRequest {
	if rs.put == nil {
		panic("put is nil")
	}

	path := rs.put.GetPath(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "PUT",
		Definition: NewServerResource(rs.put, rs.defaultContentType, rs.defaultContentType),
		Body:       body,
	}
	attachArgs(req, args)
	return req
}

func (rs *ResourceSpec) Patch(args interface{}, body interface{}) *ClientRequest {
	if rs.patch == nil {
		panic("patch is nil")
	}

	path := rs.patch.GetPath(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "PATCH",
		Definition: NewServerResource(rs.patch, rs.defaultContentType, rs.defaultContentType),
		Body:       body,
	}
	attachArgs(req, args)
	return req
}

func (rs *ResourceSpec) Delete(args interface{}) *ClientRequest {
	if rs.delete == nil {
		panic("delete is nil")
	}

	path := rs.delete.GetPath(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "DELETE",
		Definition: NewServerResource(rs.delete, rs.defaultContentType, rs.defaultContentType),
	}
	attachArgs(req, args)
	return req
}

func (rs *ResourceSpec) Head(args interface{}) *ClientRequest {
	if rs.head == nil {
		panic("head is nil")
	}

	path := rs.head.GetPath(args)
	req := &ClientRequest{
		Resource:   path,
		Verb:       "HEAD",
		Definition: NewServerResource(rs.delete, rs.defaultContentType, rs.defaultContentType),
	}
	attachArgs(req, args)
	return req
}
