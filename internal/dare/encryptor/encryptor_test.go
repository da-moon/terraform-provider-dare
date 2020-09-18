package encryptor_test

import (
	"io"

	encryptor "github.com/da-moon/terraform-provider-dare/internal/dare/encryptor"
)

func init() {
	var _ io.Writer = &encryptor.Writer{}
	var _ io.Reader = &encryptor.Reader{}
}
