package handling

import "io"

type ContentTypeDecoder struct {
	ContentType string
	Decode      func(r io.Reader, v interface{}) error
}
