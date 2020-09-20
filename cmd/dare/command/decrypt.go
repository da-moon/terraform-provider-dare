package command

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	dare "github.com/da-moon/go-dare"
	model "github.com/da-moon/go-dare/model"
	logger "github.com/da-moon/go-logger"
	primitives "github.com/da-moon/go-primitives"
	urandom "github.com/da-moon/go-urandom"
	prettyjson "github.com/hokaccha/go-prettyjson"
	cli "github.com/mitchellh/cli"
)

// DecryptCommand is a Command implementation that generates an decryption
// key.
type DecryptCommand struct {
	logFilter logger.LevelFilter
	logger    *log.Logger
	args      []string
	UI        cli.Ui
}

var _ cli.Command = &DecryptCommand{}

// Run ...
func (c *DecryptCommand) Run(args []string) int {
	const entrypoint = "decrypt"
	c.args = args
	c.UI = &cli.PrefixedUi{
		OutputPrefix: "",
		InfoPrefix:   "",
		ErrorPrefix:  "",
		Ui:           c.UI,
	}
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Info(c.Help()) }
	// [NOTE] => string array parsing must happen in the same function
	// that we are parsing cmd flags.
	var inputs []string
	cmdFlags.Var((*primitives.AppendSliceValue)(&inputs), "input", "file or directory to decrypt")
	output := OutputFlag(cmdFlags)
	masterKey := MasterKeyFlag(cmdFlags)
	masterKeyFile := MasterKeyFileFlag(cmdFlags)
	logLevel := LogLevelFlag(cmdFlags)
	nonce := NonceFlag(cmdFlags)

	err := cmdFlags.Parse(c.args)
	if err != nil {
		return 1
	}
	logGate, _, _ := c.setupLoggers(*logLevel)
	c.UI.Warn("")
	c.UI.Warn("Log data will now stream in as it occurs:\n")
	logGate.Flush()
	uuid, err := urandom.UUID()
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not generate random uuid for operation : '%v'", err))
		return 1
	}
	encKey, err := model.NewKey(c.logger, uuid,
		model.WithEncryptionKey(*masterKey),
		model.WithKeyFile(*masterKeyFile),
		model.WithNonce(*nonce),
	)
	if err != nil {
		c.UI.Error(fmt.Sprintf("[%v] %s", uuid, err.Error()))
		return 1
	}
	err = encKey.Sanitize()
	if err != nil {
		c.UI.Error(fmt.Sprintf("[%v] %s", uuid, err.Error()))
		return 1
	}
	l := logger.NewWrappedLogger(c.logger)
	l.Trace("cli-decrypt => DecryptOutputFlag '%v'", *output)
	l.Trace("cli-decrypt => MasterKeyFlag '%v'", *masterKey)
	l.Trace("cli-decrypt => MasterKeyFileFlag '%v'", *masterKeyFile)
	l.Trace("cli-decrypt => LogLevelFlag '%v'", *logLevel)
	l.Trace("cli-decrypt => NonceFlag '%v'", *nonce)
	results := make([]model.DecryptResponse, 0)
	for _, v := range inputs {
		l.Trace("cli-decrypt => InputFlag '%v'", v)
		request, err := encKey.NewDecryptRequest(v, *output)
		if err != nil {
			c.UI.Error(fmt.Sprintf("[%v] %s", uuid, err.Error()))
			return 1
		}
		result, err := dare.DecryptFile(request)
		if err != nil {
			c.UI.Error(fmt.Sprintf("[%v] %s", uuid, err.Error()))
			return 1
		}
		results = append(results, *result)
		// c.UI.Output(strings.Join(append([]string{"Decryption Result:"}, result.String()...), "\n"))
	}
	s, _ := prettyjson.Marshal(results)
	c.UI.Output(string(s))
	return 0
}

// Synopsis ...
func (c *DecryptCommand) Synopsis() string {
	return "Decrypts artifacts"
}

// Help ...
func (c *DecryptCommand) Help() string {
	helpText := `
Usage: dare decrypt [options]

  decrypts standalone files or directories at rest.

Options:

  -log-level=INFO                 log level
                                  Default: 'INFO'
  -input=/path/to/artifact        Path to decrypted artifact or artifacts 
                                  directories to decrypt.
                                  This can be specified multiple times.
                                  Default: '$PWD'
  -output=/path/to/store          Path to store decrypted artifacts.
                                  Default: same directory as the origin file
  -master-key=secret              master key used for decrypting the artifacts.
                                  master key must be a 32 byte long hex-encoded
                                  string.
                                  'DARE_MASTER_KEY' environment variable can be
                                  used for passing in this value
  -master-key-file=/path/to/key   path to a plain text file, holding master key
                                  used for decrypting the artifacts. leading
                                  and trailing whitespaces will be trimmed from the text stored in file.
                                  this takes priority over 'master-key' flag
                                  'DARE_MASTER_KEY_FILE' environment variable
                                  can be used for passing in this value.
  -nonce=secret                   random nonce used when decrypting the
                                  artifacts.
                                  nonce must be a 24 byte long hex-encoded
                                  string.
                                  'DARE_NONCE' environment variable can be
                                  used for passing in this value

`
	return strings.TrimSpace(helpText)
}

func (c *DecryptCommand) setupLoggers(logLevel string) (*logger.GatedWriter, *logger.LogWriter, io.Writer) {
	// Setup logging. First create the gated log writer, which will
	// store logs until we're ready to show them. Then create the level
	// filter, filtering logs of the specified level.
	logGate := logger.NewGatedWriter(os.Stderr)
	c.logFilter = logger.NewLevelFilter(
		logger.WithMinLevel(strings.ToUpper(logLevel)),
		logger.WithWriter(logGate),
	)
	// Create a log writer, and wrap a logOutput around it
	logWriter := logger.NewLogWriter(512)
	logOutput := io.MultiWriter(c.logFilter, logWriter)
	// Create a logger
	c.logger = log.New(logOutput, "", log.LstdFlags)
	return logGate, logWriter, logOutput
}
