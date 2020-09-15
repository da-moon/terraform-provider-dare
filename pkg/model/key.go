package model

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"strings"

	logger "github.com/da-moon/go-logger"
	stacktrace "github.com/palantir/stacktrace"
)

// Key ...
type Key struct {
	Content string                `json:"content,omitempty"`
	File    string                `json:"file,omitempty"`
	UUID    string                `json:"uuid,omitempty"`
	logger  *logger.WrappedLogger `json:"-"`
}

// Sanitize check struct fields for unacceptable values and sets defaults
func (k *Key) Sanitize(l *logger.WrappedLogger, uuid string) error {
	var err error
	if l == nil {
		err = stacktrace.NewError("key: logger was nil")
		return err
	}
	k.UUID = uuid
	k.logger = l
	// in case key content is not provided , look into file path
	// if file path is also not provided , then generate a random key
	if len(k.Content) == 0 {
		if len(k.File) != 0 {
			b, err := ioutil.ReadFile(k.File)
			if err != nil {
				k.logger.Warn("[%s] key => could not read master key file at '%s'. err: %v", uuid, k.File, err)
			}
			k.Content = strings.TrimSpace(string(b))
		}
		if len(k.Content) == 0 {
			var key [32]byte
			_, err := io.ReadFull(rand.Reader, key[:])
			if err != nil {
				err = stacktrace.Propagate(err, "could not generate random key encryption key")
				return err
			}
			k.Content = hex.EncodeToString(key[:])
			k.logger.Warn("[%s] key => encryption key was not provided. using randomly generate key '%s'", uuid, k.Content)
		}
	}
	k.logger.Trace("[%s] key => Key.Content '%s'", k.UUID, k.Content)
	k.logger.Trace("[%s] key => Key.File '%s'", k.UUID, k.File)
	return nil
}
