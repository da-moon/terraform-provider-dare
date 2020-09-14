package command

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	logger "github.com/da-moon/go-logger"
	primitives "github.com/da-moon/go-primitives"
	logutils "github.com/hashicorp/logutils"
	cli "github.com/mitchellh/cli"
)

// EncryptCommand is a Command implementation that generates an encryption
// key.
type EncryptCommand struct {
	logFilter *logutils.LevelFilter
	logger    *log.Logger
	args      []string
	Ui        cli.Ui
}

var _ cli.Command = &EncryptCommand{}

// Run ...
func (c *EncryptCommand) Run(args []string) int {
	const entrypoint = "encrypt"
	c.args = args
	c.Ui = &cli.PrefixedUi{
		OutputPrefix: "==> ",
		InfoPrefix:   "    ",
		ErrorPrefix:  "==> [ERROR]",
		Ui:           c.Ui,
	}
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Info(c.Help()) }
	var input []string
	cmdFlags.Var((*primitives.AppendSliceValue)(&input), "input", "file or directory to encrypt")
	output := EncryptOutputFlag(cmdFlags)
	masterKey := MasterKeyFlag(cmdFlags)
	masterKeyFile := MasterKeyFileFlag(cmdFlags)
	logLevel := LogLevelFlag(cmdFlags)

	err := cmdFlags.Parse(c.args)
	if err != nil {
		c.Ui.Info(c.Help())
		return 1
	}
	logGate, _, _ := c.setupLoggers(*logLevel)
	c.Ui.Info("")
	c.Ui.Output("Log data will now stream in as it occurs:\n")
	logGate.Flush()

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
		c.logger.Printf("[WARN] No master key was provided. generated random key '%s'", *masterKey)
	}
	if len(*output) == 0 {
		c.Ui.Error("output value is needed")
		c.Ui.Info(c.Help())
		return 1
	}

	c.logger.Printf("[INFO] encrypt: input path : %v", input)
	c.logger.Printf("[INFO] encrypt: output path : %s", *output)
	if len(*masterKeyFile) != 0 {
		c.logger.Printf("[INFO] encrypt: master key file: %s", *masterKeyFile)
	}
	c.logger.Printf("[INFO] encrypt: master key : %s", *masterKey)
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

func (c *EncryptCommand) setupLoggers(logLevel string) (*logger.GatedWriter, *logger.LogWriter, io.Writer) {
	// Setup logging. First create the gated log writer, which will
	// store logs until we're ready to show them. Then create the level
	// filter, filtering logs of the specified level.
	logGate := logger.NewGatedWriter(&cli.UiWriter{Ui: c.Ui})

	c.logFilter = logger.LevelFilter()
	c.logFilter.MinLevel = logutils.LogLevel(strings.ToUpper(logLevel))
	c.logFilter.Writer = logGate
	if !logger.ValidateLevelFilter(c.logFilter.MinLevel, c.logFilter) {
		c.Ui.Error(fmt.Sprintf(
			"Invalid log level: %s. Valid log levels are: %v",
			c.logFilter.MinLevel, c.logFilter.Levels))
		return nil, nil, nil
	}

	// Create a log writer, and wrap a logOutput around it
	logWriter := logger.NewLogWriter(512)
	var logOutput io.Writer
	logOutput = io.MultiWriter(c.logFilter, logWriter)
	// Create a logger
	c.logger = log.New(logOutput, "", log.LstdFlags)
	return logGate, logWriter, logOutput
}
