package handling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/gotgo/fw/tracing"
	"github.com/gotgo/fw/util"
	"github.com/gotgo/gokn/rest"
)

type ContentTypeDecoders struct {
	library map[string]*ContentTypeDecoder
}

func NewContentTypeDecoders() *ContentTypeDecoders {
	cd := &ContentTypeDecoders{
		library: make(map[string]*ContentTypeDecoder),
	}

	//add defaults
	json := &ContentTypeDecoder{
		ContentType: "application/json",
		Decode:      JsonDecoder,
	}
	cd.library[json.ContentType] = json
	return cd
}

func JsonDecoder(reader io.Reader, v interface{}, trace tracing.Tracer) error {
	if bytes, err := ioutil.ReadAll(reader); err != nil {
		return err
	} else if err = json.Unmarshal(bytes, &v); err != nil {
		return err
	} else {
		trace.Annotate(tracing.FromRequestData, "body", fmt.Sprintf("%s", bytes))
		return nil
	}
}

func (cd *ContentTypeDecoders) Get(types []string) *ContentTypeDecoder {
	for _, t := range types {
		decoder := cd.library[t]
		if decoder != nil {
			return decoder
		}
	}
	return nil
}

func (cd *ContentTypeDecoders) Set(decoder *ContentTypeDecoder) {
	cd.library[decoder.ContentType] = decoder
}

func containsType(s []string, c string) bool {
	for _, a := range s {
		if a == c {
			return true
		}
	}
	return false
}

func isBytes(t reflect.Type) bool {
	var b byte
	bt := reflect.TypeOf(b)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		return t.Elem() == bt
	default:
		return false
	}
}

func (cd *ContentTypeDecoders) DecodeBody(req *rest.Request, trace tracing.Tracer) error {
	switch req.Raw.Method {
	case "GET", "DELETE", "HEAD":
		return nil
	}

	ctype := req.Raw.Header["Content-Type"]
	decoder := cd.Get(ctype)

	//new empty instance
	body := req.Definition.RequestBody()

	var bts []byte
	var err error

	if bts, err = req.Bytes(); err != nil {
		return err
	}

	if body != nil {

		if isBytes(reflect.TypeOf(body)) {
			//if body type is castable to []byte, then we don't encode, just set directly
			req.Body = bts
			return nil
		} else if decoder != nil {
			decoder.Decode(bytes.NewReader(bts), &body, trace)
			req.Body = body
		} else if containsType(ctype, "application/x-www-form-urlencoded") {
			req.Raw.ParseForm()
			if err := util.MapHeaderToStruct(req.Raw.Form, &body); err != nil {
				return err
			}
			trace.Annotate(tracing.FromRequestData, "body", fmt.Sprintf("%+v", body))
			req.Body = body
		} else {
			if err = json.Unmarshal(bts, &body); err != nil {
				return err
			}
			trace.Annotate(tracing.FromRequestData, "body", fmt.Sprintf("%s", bts))
			req.Body = body
		}
	}
	return nil
}
