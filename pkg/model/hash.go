package model

import "github.com/palantir/stacktrace"

// Hash stores file hash values
type Hash struct {
	Path   string `json:"path,omitempty"`
	Md5    string `json:"md5,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
}

// Sanitize checks for unacceptable values
func (h *Hash) Sanitize() error {
	var err error
	if len(h.Md5) == 0 {
		err = stacktrace.NewError("returned MD5 hash of '%v' was an empty string", h.Path)
		return err
	}
	if len(h.Sha256) == 0 {
		err = stacktrace.NewError("returned Sha256 hash of '%v' was an empty string", h.Path)
		return err
	}
	return nil
}
