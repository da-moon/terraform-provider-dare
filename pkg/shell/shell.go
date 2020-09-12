package shell

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/palantir/stacktrace"
)

// Run  -
func Run(command string, args ...string) (string, error) {
	errCH := make(chan error)
	resCh := make(chan string)
	go func() {
		cmd := exec.Command(command, args...)
		// cmd.Stderr = os.Stderr
		b, err := cmd.Output()
		if err != nil {
			err = stacktrace.Propagate(err, "Error running %s %s", cmd.Path, strings.Join(cmd.Args, " "))
			errCH <- err
			return
		}
		resCh <- string(b)
		return
	}()
	select {
	case err := <-errCH:
		{
			return "", err
		}
	case res := <-resCh:
		{
			return res, nil
		}
	}
}

// RunRealtime  -
func RunRealtime(command string, args ...string) {

	cmd := exec.Command(command, args...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

type response struct {
	Results []map[string]interface{} `json:"results,omitempty"`
}

func jsonResponseAsArray(input string) ([]map[string]interface{}, error) {
	input = (strings.Join(strings.Fields((input)), " "))
	resp := &response{}
	err := json.Unmarshal([]byte(input), &resp)
	if err != nil {
		err = stacktrace.Propagate(err, "JSON unmarshalling failed with %v", input)
		return nil, err
	}

	return resp.Results, nil
}
