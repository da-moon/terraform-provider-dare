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
