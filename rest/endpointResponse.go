package rest

import (
	"io/ioutil"
	"net/http"
)

type EndpointResponse struct {
	HttpResponse *http.Response
}

func (er *EndpointResponse) Bytes() ([]byte, error) {
	defer er.HttpResponse.Body.Close()

	if contents, err := ioutil.ReadAll(er.HttpResponse.Body); err != nil {
		return nil, err
	} else {
		return contents, nil
	}
}
