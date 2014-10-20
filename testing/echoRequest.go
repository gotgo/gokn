package testing

import "net/url"

type EchoRequest struct {
	URL     *url.URL            `json:"url"`
	Headers map[string][]string `json:"header"`
	Body    interface{}         `json:"body"`
	Method  string              `json:"method"`
	Error   error               `json:"error"`
}
