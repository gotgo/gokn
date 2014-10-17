package handling

import (
	"net/http"

	"github.com/gotgo/gokn/rest"
)

type BindingFunc func(rest.HandlerFunc) func(*rest.Request, rest.Responder)

type SimpleRouter interface {
	RegisterRoute(verb, path string, f func(http.ResponseWriter, *http.Request))
	RequestArgs(req *http.Request) map[string]string
}
