package dependency

import (
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testInfoResolver(pkg string) (aur.ExtendedInfoResults, error) {
	var res []aur.InfoResult

	switch pkg {
	case "yay":
		res = []aur.InfoResult{
			{
				Depends: []string{"non-yay"},
			},
		}
	case "non-yay":
		res = []aur.InfoResult{
			{
				Depends: []string{"minecraft"},
			},
		}
	case "minecraft":
		res = []aur.InfoResult{
			{
				MakeDepends: []string{"java-18"},
			},
		}
	}

	return aur.ExtendedInfoResults{
		ResultCount: len(res),
		Results:     res,
	}, nil
}

func Test(t *testing.T) {
	exp := []string{"non-yay", "minecraft"}
	r := NewResolverWithFunction(testInfoResolver)

	it, err := r.Resolve("yay")
	assert.NoError(t, err)

	i := 0
	for it.Next() {
		assert.Equal(t, exp[i], it.Item())
		i++
	}

	assert.Equal(t, 2, i)
}
