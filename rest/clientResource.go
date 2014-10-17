package rest

type ClientResource interface {
	Get(args interface{}) *ClientRequest
	Post(args interface{}, body interface{}) *ClientRequest
	Put(args interface{}, body interface{}) *ClientRequest
	Patch(args interface{}, body interface{}) *ClientRequest
	Delete(args interface{}) *ClientRequest
	Head(args interface{}) *ClientRequest
}
