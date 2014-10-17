package rest

import (
	"fmt"

	"github.com/gotgo/fw/tracing"
)

type RequestContext struct {
	user  map[string]interface{}
	Trace tracing.Tracer
}

func NewRequestContext() *RequestContext {
	ctx := new(RequestContext)
	ctx.user = make(map[string]interface{})
	ctx.Trace = new(tracing.NopTracer)
	return ctx
}

func format(ns string, key string) string {
	return fmt.Sprintf("%s.%s", ns, key)
}

func (r *RequestContext) Add(ns string, key string, value interface{}) {
	r.user[format(ns, key)] = value
}

func (r *RequestContext) Get(ns string, key string) (value interface{}, found bool) {
	value, ok := r.user[format(ns, key)]
	return value, ok
}

func (r *RequestContext) Remove(ns string, key string) {
	r.user[format(ns, key)] = nil
}
