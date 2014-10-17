package rest

type ResourceEndpoint struct {
	Host         string `json:"host"`         // Can be host or host:port
	Scheme       string `json:"scheme"`       //http, https, websocket,
	ResourceRoot string `json:"resourceRoot"` // Root could be /api or /v1
}
