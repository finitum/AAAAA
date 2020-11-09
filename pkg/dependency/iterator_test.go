package dependency

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNext(t *testing.T) {
	deps := []Dependency{{name: "yay"}, {name: "not-yay"}, {name: "minecraft"}}
	it, err := NewIterator(deps)
	assert.NoError(t, err)

	defer it.Close()

	for i := 0; i < len(deps); i++ {
		assert.Equal(t, true, it.Next())
		assert.Equal(t, deps[i], it.Item())
	}
	assert.Equal(t, it.Next(), false)
}

func TestPush(t *testing.T) {
	deps := []Dependency{{name: "yay"}, {name: "minecraft"}}
	it, err := NewIterator(deps)
	assert.NoError(t, err)

	defer it.Close()

	assert.Equal(t, true, it.Next())
	assert.Equal(t, Dependency{name: "yay"}, it.Item())
	it.Push(Dependency{name: "yay"})
	assert.Equal(t, true, it.Next())
	assert.Equal(t, Dependency{name: "minecraft"}, it.Item())
	assert.Equal(t, true, it.Next())
	assert.Equal(t, Dependency{name: "yay"}, it.Item())
	assert.Equal(t, false, it.Next())
}
