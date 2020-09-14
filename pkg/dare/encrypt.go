package dare

import (
	"encoding/hex"
	"io"
	"os"

	config "github.com/da-moon/terraform-provider-dare/pkg/dare/config"
	encryptor "github.com/da-moon/terraform-provider-dare/pkg/dare/encryptor"
	"github.com/da-moon/terraform-provider-dare/pkg/hashsink"
	"github.com/da-moon/terraform-provider-dare/pkg/model"
	"github.com/palantir/stacktrace"
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
			return err
		}
		_, err = encWriter.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}
	}
	return nil
}

// EncryptFile encrypts a given file and store it at a given path
func EncryptFile(masterkey, input, output string) (*model.EncryptResponse, error) {
	result := &model.EncryptResponse{}
	nonce, err := RandomNonce()
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt data due to failure in generating random nonce")
		return nil, err
	}
	result.RandomNonce = hex.EncodeToString(nonce[:])
	var key [32]byte
	decoded, err := hex.DecodeString(masterkey)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt data due to failure in decoding encryption key")
		return nil, err
	}
	if len(decoded) != 32 {
		err = stacktrace.NewError("could not encrypt data since given encoded encryption key is %d bytes. We expect 32 byte keys", len(decoded))
		return nil, err
	}
	copy(key[:], decoded[:32])

	fi, err := os.Stat(input)
	if err == nil {
		if fi.Size() == 0 {
			os.Remove(input)
			err = stacktrace.NewError("decryption failure due to file with empty size at '%v'", input)
			return nil, err
		}
	}
	if err != nil {
		err = stacktrace.Propagate(err, "could not stat src at '$v'", input)
		return nil, err
	}
	srcFile, err := os.Open(input)
	if err != nil {
		err = stacktrace.Propagate(err, "could not encrypt due to failure in opening source file at %s", input)
		return nil, err
	}
	defer srcFile.Close()
	os.Remove(output)
	destinationFile, err := os.Create(output)
	if err != nil {
		err = stacktrace.NewError("could not successfully create a new empty file for %s", output)
		return nil, err
	}
	defer destinationFile.Close()
	dstWriter := hashsink.NewWriter(destinationFile)
	err = EncryptWithWriter(dstWriter, srcFile, key, nonce)
	if err != nil {
		err = stacktrace.Propagate(err, "Could not Encrypt file at '%s' and store it in '%s' ", input, output)
		return nil, err
	}
	result.OutputHash = &model.Hash{
		Path:   output,
		Md5:    dstWriter.MD5HexString(),
		Sha256: dstWriter.SHA256HexString(),
	}
	return result, nil
}
