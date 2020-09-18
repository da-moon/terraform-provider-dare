package command

import (
	"flag"
	"os"

	"github.com/da-moon/go-primitives"
)

// DDSizeFlag ...
func DDSizeFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_DEMO_SIZE")
	return f.String("size", result,
		"demo file size.")
}

// DDPathFlag ...
func DDPathFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_DEMO_PATH")
	return f.String("path", result,
		"path to store demo file.")
}

// MasterKeyFlag ...
func MasterKeyFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_MASTER_KEY")
	return f.String("master-key", result,
		"Master Key used in encryption-decryption process.")
}

// MasterKeyFileFlag ...
func MasterKeyFileFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_MASTER_KEY_FILE")
	return f.String("master-key-file", result,
		"plain text file holding Master Key used in encryption-decryption process.")
}

// EncryptOutputFlag ...
func EncryptOutputFlag(f *flag.FlagSet) *string {
	var result string
	return f.String("output", result,
		"Path to store encrypted artifacts.")
}

// DecryptOutputFlag ...
func DecryptOutputFlag(f *flag.FlagSet) *string {
	dir, _ := os.Getwd()
	result := primitives.PathJoin(dir, "decrypted")
	return f.String("output", result,
		"Path to store decrypted artifacts.")
}

// LogLevelFlag ...
func LogLevelFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_LOG_LEVEL")
	if result == "" {
		result = "INFO"
	}
	return f.String("log-level", result,
		"flag used to indicate log level")
}

// RegexFlag ...
func RegexFlag(f *flag.FlagSet) *string {
	var result string
	return f.String("regex", result,
		"regex used for recursive file search")
}

// NonceFlag ...
func NonceFlag(f *flag.FlagSet) *string {
	result := os.Getenv("DARE_NONCE")
	return f.String("nonce", result,
		"random initial nonce used when encrypting artifacts")
}
