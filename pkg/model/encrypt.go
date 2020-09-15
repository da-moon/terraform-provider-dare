package model

import (
	"log"
	"path/filepath"

	logger "github.com/da-moon/go-logger"
	"github.com/da-moon/go-primitives"
	stacktrace "github.com/palantir/stacktrace"
)

// EncryptRequest ...
type EncryptRequest struct {
	Key    Key                   `json:"key,omitempty"`
	Files  []File                `json:"files,omitempty"`
	UUID   string                `json:"uuid,omitempty"`
	logger *logger.WrappedLogger `json:"-"`
}

// Sanitize ...
func (r *EncryptRequest) Sanitize(l *log.Logger, uuid string) error {
	var err error
	if l == nil {
		err = stacktrace.NewError("logger was nil")
		return err
	}
	r.UUID = uuid
	r.logger = logger.NewWrappedLogger(l)
	err = r.Key.Sanitize(r.logger, uuid)
	if err != nil {
		err = stacktrace.Propagate(err, "encrypt-request: could not sanitize keys", uuid)
		return err
	}

	if len(r.Files) == 0 {
		err = stacktrace.NewError("encrypt-request: files to process list is empty", uuid)
		return err
	}
	for _, v := range r.Files {
		err = v.Sanitize(r.logger, uuid)
		if err != nil {
			err = stacktrace.Propagate(err, "encrypt-request: could not sanitize files list", uuid)
			return err
		}
	}
	return nil
}

// AddTarget adds a target file to encrypt
func (r *EncryptRequest) AddTarget(source, destination, regex string) {
	if r.Files == nil {
		r.Files = make([]File, 0)
	}
	f, fi, err := primitives.OpenPath(source)
	if err != nil {
		r.logger.Error("[%s] encrypt-request => could not open '%s' : '%v'", r.UUID, source, err)
		return
	}
	defer f.Close()
	if !fi.IsDir() {
		if filepath.Ext(source) != "enc" {
			parent := filepath.Dir(source)
			if len(destination) != 0 {
				parent = primitives.PathJoin(destination, parent)
			}
			r.Files = append(r.Files, File{
				Source:      source,
				Destination: primitives.PathJoin(parent, filepath.Base(source)+".enc"),
			})
		}
	} else {
		files, err := primitives.FindFile(source, regex)
		if err != nil {
			r.logger.Error("[%s] encrypt-request => could search '%s' for '%s' pattern : '%v'", r.UUID, source, regex, err)
			return
		}
		for _, v := range files {
			if filepath.Ext(source) != "enc" {
				parent := filepath.Dir(v)
				if len(destination) != 0 {
					parent = primitives.PathJoin(destination, parent)
				}
				r.Files = append(r.Files, File{
					Source:      v,
					Destination: primitives.PathJoin(parent, filepath.Base(v)+".enc"),
				})
			}
		}
	}
}

// EncryptResponse ...
type EncryptResponse struct {
	OutputHash  *Hash  `json:"output_hash,omitempty"`
	RandomNonce string `json:"random_nonce,omitempty"`
	RandomKey   string `json:"random_key,omitempty"`
}

// func (r *EncryptResponse) Sanitize(logger *log.Logger) error {
// 	var err error
// 	if len(r.Source) == 0 {
// 		err = stacktrace.NewError("input path for ")
// 	}
// }
