package handling

type ContentTypeEncoder struct {
	ContentType string
	Encode      func(v interface{}) ([]byte, error)
}
