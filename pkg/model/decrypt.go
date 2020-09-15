package model

// DecryptRequest ...
type DecryptRequest struct {
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	Nonce       string `json:"nonce,omitempty"`
	Key         string `json:"key,omitempty"`
}
