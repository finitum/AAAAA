package cmdutil

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func RunCommand(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
	cmd.Stderr = io.MultiWriter(&stderr, os.Stdout)

	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}
