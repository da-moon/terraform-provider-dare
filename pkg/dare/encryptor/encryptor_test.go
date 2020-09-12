package encryptor_test

import (
	"io"

	encryptor "github.com/da-moon/terraform-provider-dare/pkg/dare/encryptor"
)

func init() {
	var _ io.Writer = &encryptor.Writer{}
	var _ io.Reader = &encryptor.Reader{}
}
