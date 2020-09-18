package command

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

// KeygenCommand is a Command implementation that generates an encryption
// key.
type KeygenCommand struct {
	UI cli.Ui
}

var _ cli.Command = &KeygenCommand{}

// Run ...
func (c *KeygenCommand) Run(_ []string) int {
	const length = 32
	key := make([]byte, length)
	n, err := rand.Reader.Read(key)
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not read random data: %s", err))
		return 1
	}
	if n != length {
		c.UI.Error("could not read enough entropy. Generate more entropy!")
		return 1
	}
	c.UI.Output(hex.EncodeToString(key))
	return 0
}

// Synopsis ...
func (c *KeygenCommand) Synopsis() string {
	return "Generates a new encryption key"
}

// Help ...
func (c *KeygenCommand) Help() string {
	helpText := `
Usage: dare keygen

  Generates a new 32 byte long hex encoded encryption key that can be used to for encrypting data.
`
	return strings.TrimSpace(helpText)
}
