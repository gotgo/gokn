package rest

import "encoding/json"

type FrozenRequest struct {
	Headers  map[string][]string `json:"headers"`
	Verb     string              `json:"verb"`
	Resource string              `json:"resource"`
	Body     json.RawMessage     `json:"body"`

	ResourceT string        `json:"resourceT,omitempty"`
	Requestor *FrozenSender `json:"requestor,omitempty"`
}

type FrozenSender struct {
	Device   string `json:"device,omitempty"`
	User     string `json:"user,omitempty"`
	Account  string `json:"account,omitempty"`
	Location string `json:"location,omitempty"`
}

type FrozenResponse struct {
	Headers    map[string][]string `json:"headers"`
	StatusCode string              `json:"statusCode"`
	Status     string              `json:"status"`
	Body       json.RawMessage     `json:"body"`
}

type FrozenCommunication struct {
	Request  *FrozenRequest  `json:"request"`
	Response *FrozenResponse `json:"response"`
}
