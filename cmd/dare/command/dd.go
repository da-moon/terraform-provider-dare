package command

import (
	"flag"
	"fmt"
	"math"
	mathrand "math/rand"
	"os"
	"strings"

	codec "github.com/da-moon/go-codec"
	model "github.com/da-moon/go-dare/model"
	primitives "github.com/da-moon/go-primitives"
	stream "github.com/da-moon/go-stream"
	cli "github.com/mitchellh/cli"
	stacktrace "github.com/palantir/stacktrace"
)

// DDCommand is a Command implementation that generates an encryption
// key.
type DDCommand struct {
	args []string
	UI   cli.Ui
}

var _ cli.Command = &DDCommand{}

// Run ...
func (c *DDCommand) Run(args []string) int {
	c.UI = &cli.PrefixedUi{
		OutputPrefix: "==> ",
		Ui:           c.UI,
	}

	c.args = args
	const entrypoint = "dd"
	cmdFlags := flag.NewFlagSet(entrypoint, flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Info(c.Help()) }
	sizeString := DDSizeFlag(cmdFlags)
	pathString := DDPathFlag(cmdFlags)
	err := cmdFlags.Parse(c.args)
	if err != nil {
		c.UI.Info(c.Help())
		return 1
	}
	if len(*sizeString) == 0 {
		c.UI.Error("size value is needed")
		c.UI.Info(c.Help())
		return 1
	}
	if len(*pathString) == 0 {
		c.UI.Error("path value is needed")
		c.UI.Info(c.Help())
		return 1
	}
	parsedSize, err := primitives.FileSizeStringToInt(*sizeString)
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not parse given size: %s", err))
		c.UI.Info(c.Help())
		return 1
	}
	os.Remove(*pathString)
	result, err := createRandomFile(*pathString, int(parsedSize))
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not create random file: %s", err))
		c.UI.Info(c.Help())
		return 1
	}
	if result == nil {
		c.UI.Error("could not create random file")
		c.UI.Info(c.Help())
		return 1
	}
	err = result.Sanitize()
	if err != nil {
		c.UI.Error(fmt.Sprintf("there were an issue with random file metadata: %s", err))
		c.UI.Info(c.Help())
		return 1
	}
	c.UI.Output(fmt.Sprintf("output path : %s", *pathString))
	c.UI.Output(fmt.Sprintf("MD5 Hash : %s", result.Md5))
	c.UI.Output(fmt.Sprintf("SHA256 Hash : %s", result.Sha256))
	return 0
}

// Synopsis ...
func (c *DDCommand) Synopsis() string {
	return "Generates a new file used for testing"
}

// Help ...
func (c *DDCommand) Help() string {
	helpText := `
Usage: dare dd [options]

  generates a new human readable JSON lorem ipsum file. 

Options:

  --size=1MB file size to generate.
  --path=/tmp/plain target path to store the file.
`
	return strings.TrimSpace(helpText)
}

func createRandomFile(path string, maxSize int) (*model.Hash, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		err = stacktrace.Propagate(err, "Can't open %s for writing", path)
		return nil, err
	}
	defer file.Close()
	hashWriter, err := stream.NewHashWriter(
		file,
		stream.WithMD5(),
		stream.WithSHA256(),
	)
	if err != nil {
		err = stacktrace.Propagate(err, "can't get a hashwriter for writing to %s", path)
		return nil, err
	}
	size := maxSize/2 + mathrand.Int()%(maxSize/2)
	loremString := path + `---Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin facilisis mi sapien, vitae accumsan libero malesuada in. Suspendisse sodales finibus sagittis. Proin et augue vitae dui scelerisque imperdiet. Suspendisse et pulvinar libero. Vestibulum id porttitor augue. Vivamus lobortis lacus et libero ultricies accumsan. Donec non feugiat enim, nec tempus nunc. Mauris rutrum, diam euismod elementum ultricies, purus tellus faucibus augue, sit amet tristique diam purus eu arcu. Integer elementum urna non justo fringilla fermentum. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Quisque sollicitudin elit in metus imperdiet, et gravida tortor hendrerit. In volutpat tellus quis sapien rutrum, sit amet cursus augue ultricies. Morbi tincidunt arcu id commodo mollis. Aliquam laoreet purus sed justo pulvinar, quis porta risus lobortis. In commodo leo id porta mattis.`
	byteSizeOfDefaultLorem := len([]byte(loremString))
	repetitions := int(math.Round(float64(size / byteSizeOfDefaultLorem)))
	for i := 0; i < repetitions; i++ {
		enc, _ := codec.EncodeJSONWithIndentation(map[int]string{
			i: (loremString),
		})
		hashWriter.Write([]byte(enc))
	}
	md5Hex, err := hashWriter.HexString(stream.MD5)
	if err != nil {
		err = stacktrace.Propagate(err, "could not calculate hex encoded md5 hash of '%s' ", path)
		return nil, err
	}
	sha256Hex, err := hashWriter.HexString(stream.SHA256)
	if err != nil {
		err = stacktrace.Propagate(err, "could not calculate hex encoded sha256 hash of '%s' ", path)
		return nil, err
	}
	result := &model.Hash{
		Md5:    md5Hex,
		Sha256: sha256Hex,
	}
	return result, nil
}
