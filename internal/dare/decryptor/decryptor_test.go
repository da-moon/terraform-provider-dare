package decryptor_test

import (
	"io"

	decryptor "github.com/da-moon/terraform-provider-dare/internal/dare/decryptor"
)

func init() {
	var _ io.Writer = &decryptor.Writer{}
	var _ io.Reader = &decryptor.Reader{}
}
