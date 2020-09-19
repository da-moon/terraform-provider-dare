package dare

import (
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	files "github.com/da-moon/go-files"
	config "github.com/da-moon/terraform-provider-dare/internal/dare/config"
	encryptor "github.com/da-moon/terraform-provider-dare/internal/dare/encryptor"
	hashsink "github.com/da-moon/terraform-provider-dare/pkg/hashsink"
	model "github.com/da-moon/terraform-provider-dare/pkg/model"
	stacktrace "github.com/palantir/stacktrace"
)

// EncryptWithWriter encrypts data with a passed key as it is writing it
// to an io stream (eg socket , file).
func EncryptWithWriter(
	dstwriter io.Writer,
	srcReader io.Reader,
	key [32]byte,
	nonce [24]byte,
) error {
	encWriter := encryptor.NewWriter(dstwriter, nonce, &key)
	for {
		buffer := make([]byte, config.DefaultChunkSize)
		bytesRead, err := srcReader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			err = stacktrace.Propagate(err, "could not read from source")
			return err
		}
		_, err = encWriter.Write(buffer[:bytesRead])
		if err != nil {
			err = stacktrace.Propagate(err, "could not write to destination")
			return err
		}
	}
	return nil
}

// EncryptFile encrypts a given file and store it at a given path
func EncryptFile(r *model.EncryptRequest) (*model.EncryptResponse, error) {
	var err error
	if r == nil {
		err = stacktrace.NewError("nil request")
		return nil, err
	}
	err = r.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not sanitize request")
		return nil, err
	}
	result := r.Response()
	result.RandomNonce = hex.EncodeToString(r.Nonce[:])
	for k, v := range r.Targets {
		srcFile, _, err := files.SafeOpenPath(k)
		if err != nil {
			err = stacktrace.Propagate(err, "could not encrypt '%s'", k)
			return nil, err
		}
		defer srcFile.Close()

		os.Remove(v)
		files.MkdirAll(filepath.Dir(v))
		destinationFile, err := os.Create(v)
		if err != nil {
			err = stacktrace.NewError("could not successfully create a new empty file for %s", v)
			return nil, err
		}
		defer destinationFile.Close()
		dstWriter := hashsink.NewWriter(destinationFile)
		err = EncryptWithWriter(dstWriter, srcFile, r.Key, r.Nonce)
		if err != nil {
			err = stacktrace.Propagate(err, "Could not Encrypt file at '%s' and store it in '%s' ", k, v)
			return nil, err
		}
		result.EncryptedArtifacts[v] = model.Hash{
			Md5:    dstWriter.MD5HexString(),
			Sha256: dstWriter.SHA256HexString(),
		}
	}
	err = result.Sanitize()
	if err != nil {
		err = stacktrace.Propagate(err, "could not Sanitize Encrypt Response")
		return nil, err
	}
	return result, nil
}
