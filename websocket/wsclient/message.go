package wsclient

type message struct {
	ID    string `json:"id,omitempty"`
	Act   string `json:"act,omitempty"`
	Data  any    `json:"data,omitempty"`
	Error *error `json:"error,omitempty"`
	Reply any    `json:"reply,omitempty"`
}
