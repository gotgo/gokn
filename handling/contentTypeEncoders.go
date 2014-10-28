package handling

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type ContentTypeEncoders struct {
	library map[string]*ContentTypeEncoder
}

func NewContentTypeEncoders() *ContentTypeEncoders {
	cd := &ContentTypeEncoders{
		library: make(map[string]*ContentTypeEncoder),
	}

	json := &ContentTypeEncoder{
		ContentType: "application/json",
		Encode:      jsonEncoder,
	}

	text := &ContentTypeEncoder{
		ContentType: "text/plain",
		Encode:      plainTextEncoder,
	}

	cd.library[json.ContentType] = json
	cd.library[text.ContentType] = text
	return cd
}

func jsonEncoder(v interface{}) ([]byte, error) {
	if bytes, err := json.Marshal(&v); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func plainTextEncoder(v interface{}) ([]byte, error) {
	if str, ok := v.(string); !ok {
		return nil, errors.New("failed to cast to string")
	} else {
		return []byte(str), nil
	}
}

func (cte *ContentTypeEncoders) Set(encoder *ContentTypeEncoder) {
	cte.library[encoder.ContentType] = encoder
}

func (cte *ContentTypeEncoders) Encode(data interface{}, contentType string) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	} else if bts, ok := data.([]byte); ok {
		return bts, nil //bytes pass through
	} else if rdr, ok := data.(io.Reader); ok {
		if bts, err := ioutil.ReadAll(rdr); err != nil {
			return nil, err
		} else {
			return bts, nil
		}
	}

	if encoder := cte.library[contentType]; encoder != nil {
		return encoder.Encode(data)
	}

	return nil, errors.New("Encode Fail.  Unknown contentType " + contentType)
}
