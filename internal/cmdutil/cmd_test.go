package cmdutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommandStdout(t *testing.T) {
	stdout, _, err := RunCommand("/bin/bash", "-c", "builtin echo hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello\n", stdout)
	// assert.Empty(t, stderr)
}

func TestRunCommandStdErr(t *testing.T) {
	stdout, stderr, err := RunCommand("/bin/bash", "-c", "builtin echo hello 1>&2")
	assert.NoError(t, err)
	assert.Empty(t, stdout)
	assert.Equal(t, "hello\n", stderr)
}

func TestRunCommandNolinebreak(t *testing.T) {
	stdout, stderr, err := RunCommand("/bin/bash", "-c", "builtin echo -n hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", stdout)
	assert.Empty(t, stderr)
	fmt.Print("\n")
}
