package model

type EncryptResponse struct {
	OutputHash  *Hash  `json:"output_hash,omitempty"`
	RandomNonce string `json:"random_nonce,omitempty"`
	RandomKey   string `json:"random_key,omitempty"`
}
