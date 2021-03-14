package dependency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	deps := []Dependency{{Name: "yay"}, {Name: "not-yay"}, {Name: "minecraft"}}
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
	deps := []Dependency{{Name: "yay"}, {Name: "minecraft"}}
	it, err := NewIterator(deps)
	assert.NoError(t, err)

	defer it.Close()

	assert.Equal(t, true, it.Next())
	assert.Equal(t, Dependency{Name: "yay"}, it.Item())
	it.Push(Dependency{Name: "yay"})
	assert.Equal(t, true, it.Next())
	assert.Equal(t, Dependency{Name: "minecraft"}, it.Item())
	assert.Equal(t, true, it.Next())
	assert.Equal(t, Dependency{Name: "yay"}, it.Item())
	assert.Equal(t, false, it.Next())
}
