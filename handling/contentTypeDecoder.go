package handling

import (
	"io"

	"github.com/gotgo/fw/tracing"
)

type ContentTypeDecoder struct {
	ContentType string
	Decode      func(r io.Reader, v interface{}, trace tracing.Tracer) error
}
