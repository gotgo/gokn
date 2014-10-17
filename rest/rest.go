package rest

import (
	"errors"
	"net/http"
)

//more work to support http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

func GetHeaderValue(key string, headers map[string][]string) string {
	values := headers[key]
	if values != nil && len(values) > 0 {
		return values[0]
	}
	return ""
}

// Bytes, is a helper method to reduce the number of lines to get a byte array out of the
// EndpointResponse
func Bytes(resp *EndpointResponse, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if resp.HttpResponse == nil {
		return nil, errors.New("No HttpResponse Response")
	}
	if resp.HttpResponse.StatusCode != http.StatusOK {
		return nil, errors.New(resp.HttpResponse.Status)
	}

	if bytes, err := resp.Bytes(); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}
