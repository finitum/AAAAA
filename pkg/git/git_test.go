package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatestSHA(t *testing.T) {
	url := "https://github.com/finitum/AAAAA"

	hash, err := LatestHash(url, "main")
	assert.NoError(t, err)

	fmt.Println(hash.String())
}
