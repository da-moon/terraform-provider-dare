package dare

import (
	"io"
	"os"
	"path/filepath"

	"github.com/da-moon/go-files"
	config "github.com/da-moon/terraform-provider-dare/internal/dare/config"
	decryptor "github.com/da-moon/terraform-provider-dare/internal/dare/decryptor"
	hashsink "github.com/da-moon/terraform-provider-dare/pkg/hashsink"
	model "github.com/da-moon/terraform-provider-dare/pkg/model"
	stacktrace "github.com/palantir/stacktrace"
)

// DecryptWithWriter ...
func DecryptWithWriter(
	dstwriter io.Writer,
	srcReader io.Reader,
	key [32]byte,
	nonce [24]byte,
) error {
	decWriter := decryptor.NewWriter(dstwriter, nonce, &key)
	for {
		buffer := make([]byte, config.DefaultChunkSize+config.DefaultOverhead)
		bytesRead, err := srcReader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		_, err = decWriter.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptFile ...
func DecryptFile(r *model.DecryptRequest) (*model.DecryptResponse, error) {
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
	for k, v := range r.Targets {
		srcFile, _, err := files.SafeOpenPath(k)
		if err != nil {
			err = stacktrace.Propagate(err, "could not decrypt '%s'", k)
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
			err = stacktrace.Propagate(err, "Could not decrypt file at '%s' and store it in '%s' ", k, v)
			return nil, err
		}
		result.DecryptedArtifacts[v] = model.Hash{
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
