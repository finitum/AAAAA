package dependency

import (
	"errors"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
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
	case "circ-a":
		res = []aur.InfoResult{
			{
				Depends:     []string{"circ-b"},
				MakeDepends: []string{"minecraft"},
			},
		}
	case "circ-b":
		res = []aur.InfoResult{
			{
				Depends:     []string{"circ-a", "non-yay"},
				MakeDepends: []string{"minecraft"},
			},
		}
	}

	return aur.ExtendedInfoResults{
		ResultCount: len(res),
		Results:     res,
	}, nil
}

func TestNormal(t *testing.T) {
	exp := make(map[string]Dependency)
	exp["non-yay"] = Dependency{name: "non-yay", dependencies: []string{"minecraft"}}
	exp["minecraft"] = Dependency{name: "minecraft", dependencies: []string{}}

	actual := make([]Dependency, 0, len(exp))
	r := NewResolverWithFunction(testInfoResolver)

	it, err := r.Resolve("yay")
	assert.NoError(t, err)

	for it.Next() {
		actual = append(actual, it.Item())
	}

	verify(exp, actual, t)
}

func TestAurErr(t *testing.T) {
	r := NewResolverWithFunction(func(_ string) (aur.ExtendedInfoResults, error) {
		return aur.ExtendedInfoResults{}, errors.New("err")
	})

	_, err := r.Resolve("something")
	assert.Error(t, err)
}

func TestAurToManyResults(t *testing.T) {
	r := NewResolverWithFunction(func(_ string) (aur.ExtendedInfoResults, error) {
		return aur.ExtendedInfoResults{ResultCount: 42}, nil
	})

	_, err := r.Resolve("something")
	assert.Error(t, err)
}

func TestCircular(t *testing.T) {
	timeout := time.After(2 * time.Second)
	done := make(chan struct{})

	go func() {
		exp := make(map[string]Dependency)
		exp["circ-b"] = Dependency{name: "circ-b", dependencies: []string{"circ-a", "non-yay", "minecraft"}}
		exp["non-yay"] = Dependency{name: "non-yay", dependencies: []string{"minecraft"}}
		exp["minecraft"] = Dependency{name: "minecraft", dependencies: []string{}}

		actual := make([]Dependency, 0, len(exp))
		r := NewResolverWithFunction(testInfoResolver)

		it, err := r.Resolve("circ-a")
		assert.NoError(t, err)

		for it.Next() {
			actual = append(actual, it.Item())
		}

		verify(exp, actual, t)

		done <- struct{}{}
	}()

	select {
	case <-timeout:
		t.Fatal("dependency resolver is probably stuck in a loop, while it should handle circular dependencies")
	case <-done:
	}
}

func verify(exp map[string]Dependency, actual []Dependency, t *testing.T) {
	assert.Equal(t, len(exp), len(actual))

	for _, dep := range actual {
		eDep, ok := exp[dep.name]
		assert.True(t, ok)

		sort.Strings(eDep.dependencies)
		sort.Strings(dep.dependencies)
		assert.EqualValues(t, eDep.dependencies, dep.dependencies)
	}
}
