package git

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLatestSHA(t *testing.T) {
	url := "https://github.com/finitum/AAAAA"

	hash, err := LatestHash(url, "main")
	assert.NoError(t, err)

	fmt.Println(hash.String())
}
