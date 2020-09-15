package dare

import (
	"io"

	config "github.com/da-moon/terraform-provider-dare/pkg/dare/config"
	decryptor "github.com/da-moon/terraform-provider-dare/pkg/dare/decryptor"
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

// func DecryptFile(req *model.DecryptRequest) (*model.DecryptResponse, error) {
// 	result := &model.DecryptResponse{}
// 	decKey, err := hex.DecodeString(req.Key)
// 	if err != nil {
// 		err = stacktrace.Propagate(err, "could not decrypt file '%s' due to failure in decoding encryption key from hex string '%s'", req.Source, req.Key)
// 		return nil, err
// 	}
// 	if len(decKey) != 32 {
// 		err = stacktrace.NewError("could not decrypt data since given encoded encryption key is %d bytes. We expect 32 byte keys", len(decKey))
// 		return nil, err

// 	}
// 	var key [32]byte
// 	copy(key[:], decKey[:32])
// 	decNonce, err := hex.DecodeString(req.Nonce)
// 	if err != nil {
// 		err = stacktrace.Propagate(err, "could not decrypt file '%s' due to failure in decoding nonce from hex string '%s'", req.Source, req.Nonce)
// 		return nil, err
// 	}
// 	if len(decNonce) != 24 {
// 		err = stacktrace.NewError("could not decrypt data since given encoded nonce is %d bytes. We expect 24 byte keys", len(decNonce))
// 		return nil, err

// 	}
// 	var nonce [24]byte
// 	copy(nonce[:], decNonce[:24])
// 	fi, err := os.Stat(req.Source)
// 	if err == nil {
// 		if fi.Size() == 0 {
// 			os.Remove(req.Source)
// 			err = stacktrace.NewError("decryption failure due to file with empty size at '%v'", req.Source)
// 			return nil, err
// 		}
// 	}
// 	if err != nil {
// 		err = stacktrace.Propagate(err, "could not stat src at '$v'", req.Source)
// 		return nil, err
// 	}
// 	encSourceFile, err := os.Open(req.Source)
// 	if err != nil {
// 		err = stacktrace.Propagate(err, "could not decrypt due to failure in opening source file at %s", req.Source)
// 		return nil, err
// 	}
// 	defer encSourceFile.Close()
// 	os.Remove(req.Destination)
// 	decDestinationFile, err := os.Create(req.Destination)
// 	if err != nil {
// 		err = stacktrace.NewError("could not successfully create a new empty file for %s", req.Destination)
// 		return nil, err
// 	}
// 	defer decDestinationFile.Close()
// 	dstWriter := hashsink.NewWriter(decDestinationFile)
// 	err = dare.DecryptWithWriter(dstWriter, encSourceFile, key, nonce)
// 	if err != nil {
// 		err = stacktrace.Propagate(err, "Could not Decrypt file at '%s' and store it in '%s' with key '%s' and nonce '%s'", req.Source, req.Destination, req.Key, req.Nonce)
// 		return nil, err
// 	}
// 	result.OutputHash = &model.Hash{
// 		Md5:    dstWriter.MD5HexString(),
// 		Sha256: dstWriter.SHA256HexString(),
// 	}
// 	return result, nil
// }
