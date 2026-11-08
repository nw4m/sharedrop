package cmd

import (
	"bytes"
	"os/exec"
)

// execCommand runs an external command and returns its stdout as string
func execCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
