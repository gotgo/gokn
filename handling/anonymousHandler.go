package handling

import "github.com/krave-n/go/rest"

func AnonymousHandler(handler rest.HandlerFunc) func(*rest.Request, rest.Responder) {
	return handler
}
