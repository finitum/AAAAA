package cmdutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCommandStdout(t *testing.T) {
	stdout, _, err := RunCommand("/usr/bin/echo", "hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
	//assert.Empty(t, stderr)
}

func TestRunCommandStdErr(t *testing.T) {
	_, stderr, err := RunCommand("/usr/bin/bash", "-c", "echo hello 1>&2")
	assert.NoError(t, err)
	//assert.Empty(t, stdout)
	assert.Equal(t, "hello\n", stderr)
}

func TestRunCommandNolinebreak(t *testing.T) {
	stdout, _, err := RunCommand("/usr/bin/echo", "-n", "hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", stdout)
	//assert.Empty(t, stderr)
	fmt.Print("\n")
}
