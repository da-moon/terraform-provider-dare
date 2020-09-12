package command

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	primitives "github.com/da-moon/go-primitives"
	dare "github.com/da-moon/terraform-provider-dare/pkg/dare"
	hashsink "github.com/da-moon/terraform-provider-dare/pkg/hashsink"
	model "github.com/da-moon/terraform-provider-dare/pkg/model"
	cli "github.com/mitchellh/cli"
	stacktrace "github.com/palantir/stacktrace"
)

// EncryptCommand is a Command implementation that generates an encryption
// key.
type EncryptCommand struct {
	args []string
	Ui   cli.Ui
}

var _ cli.Command = &EncryptCommand{}

// Run ...
func (c *EncryptCommand) Run(args []string) int {
	c.args = args
	const entrypoint = "encrypt"
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Info(c.Help()) }
	var input []string
	cmdFlags.Var((*primitives.AppendSliceValue)(&input), "input", "file or directory to encrypt")
	output := EncryptOutputFlag(cmdFlags)
	masterKey := MasterKeyFlag(cmdFlags)
	masterKeyFile := MasterKeyFileFlag(cmdFlags)
	err := cmdFlags.Parse(c.args)
	if err != nil {
		c.Ui.Info(c.Help())
		return 1
	}
	if len(input) == 0 {
		c.Ui.Error("input value is needed")
		c.Ui.Info(c.Help())
		return 1
	}

	if len(*masterKeyFile) != 0 {
		b, err := ioutil.ReadFile(*masterKeyFile)
		if err != nil {
			c.Ui.Warn((fmt.Sprintf("could not read master key file at %s : %v", *masterKeyFile, err)))
		}
		*masterKey = strings.TrimSpace(string(b))
	}
	if len(*masterKey) == 0 {
		var key [32]byte
		_, err := io.ReadFull(rand.Reader, key[:])
		if err != nil {
			c.Ui.Error(fmt.Sprintf("could not generate random key encryption key : %v", err))
			return 1
		}
		*masterKey = hex.EncodeToString(key[:])
		c.Ui.Warn(fmt.Sprintf("No master key was provided. generated random key '%s'", *masterKey, err))
	}
	if len(*output) == 0 {
		c.Ui.Error("output value is needed")
		c.Ui.Info(c.Help())
		return 1
	}
	c.Ui.Output(fmt.Sprintf("input path : %v", input))
	c.Ui.Output(fmt.Sprintf("output path : %s", *output))
	c.Ui.Output(fmt.Sprintf("master key file: %s", *masterKeyFile))
	c.Ui.Output(fmt.Sprintf("master key : %s", *masterKey))
	return 0
}

// Synopsis ...
func (c *EncryptCommand) Synopsis() string {
	return "Encrypts artifacts"
}

// Help ...
func (c *EncryptCommand) Help() string {
	helpText := `
Usage: dare encrypt [options]

  encrypts standalone files or directories at rest.

Options:

  -input=/path/to/artifact        Path to artifact or artifacts directories to
                                  encrypt.
                                  This can be specified multiple times.
  -output=/path/to/store          Path to store encrypted artifacts.
                                  Default: '$PWD/encrypted'
  -master-key=secret              master key used for encrypting the artifacts.
                                  master key must be a 32 byte long hex-encoded
                                  string.
                                  'DARE_MASTER_KEY' environment variable can be
                                  used for passing in this value
  -master-key-file=/path/to/key   path to a plain text file, holding master key
                                  used for encrypting the artifacts. leading
                                  and trailing whitespaces will be trimmed from the text stored in file.
                                  this takes priority over 'master-key' flag
                                  'DARE_MASTER_KEY_FILE' environment variable
                                  can be used for passing in this value.
`
	return strings.TrimSpace(helpText)
}

// Encrypt - Implementation of Encrypt method for go engine
func encrypt(masterkey, input, output string) (*model.EncryptResponse, error) {
	result := &model.EncryptResponse{}
	nonce, err := dare.RandomNonce()
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
	err = dare.EncryptWithWriter(dstWriter, srcFile, key, nonce)
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
