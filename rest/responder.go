package rest

type Responder interface {
	AddHeader(key, value string)
	SetBody(interface{})
	SetContentType(ct string)
	SetStatus(statusCode int, statusMessage string)
}
