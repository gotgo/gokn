package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/oleiade/reflections"
)

//more work to support http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

func GetHeaderValue(key string, headers map[string][]string) string {
	values := headers[key]
	if values != nil && len(values) > 0 {
		return values[0]
	}
	return ""
}

func FullApiWithHandlers(apiSpec interface{}) map[ServerResource]Handler {
	handlers := make(map[ServerResource]Handler)

	items, err := reflections.Items(apiSpec)
	if err != nil {
		panic(err)
	}

	for _, v := range items {
		if target, ok := v.(*ResourceSpec); ok {
			verbs, handler := target.ServeAll()
			if handler == nil {
				panic(fmt.Sprintf("no handler for endpoint %v ", target))
			}
			for _, verb := range verbs {
				handlers[verb] = handler
			}
		}
	}
	return handlers
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
