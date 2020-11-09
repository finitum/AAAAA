package dependency

import (
	"errors"
	"fmt"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func testInfoResolver(_, pkg string) (aur.InfoResult, error) {
	var res aur.InfoResult

	switch pkg {
	case "yay":
		res.Depends = []string{"non-yay"}
	case "non-yay":
		res.Depends = []string{"minecraft"}
	case "minecraft":
		res.MakeDepends = []string{"java-18"}
	case "circ-a":
		res.Depends = []string{"circ-b"}
		res.MakeDepends = []string{"minecraft"}
	case "circ-b":
		res.Depends = []string{"circ-a", "non-yay"}
		res.MakeDepends = []string{"minecraft"}
	default:
		return aur.InfoResult{}, aur.NotInAurErr
	}

	return res, nil
}

func TestNormal(t *testing.T) {
	exp := make(map[string]Dependency)
	exp["non-yay"] = Dependency{Name: "non-yay", Dependencies: []string{"minecraft"}}
	exp["minecraft"] = Dependency{Name: "minecraft", Dependencies: []string{}}

	actual := make([]Dependency, 0, len(exp))
	r := NewResolverWithFunction("", testInfoResolver)

	it, err := r.Resolve("yay")
	assert.NoError(t, err)

	for it.Next() {
		actual = append(actual, it.Item())
	}

	verify(exp, actual, t)
}

func TestAurErr(t *testing.T) {
	r := NewResolverWithFunction("", func(_, _ string) (aur.InfoResult, error) {
		return aur.InfoResult{}, errors.New("err")
	})

	_, err := r.Resolve("something")
	assert.Error(t, err)
}

func TestCircular(t *testing.T) {
	timeout := time.After(2 * time.Second)
	done := make(chan struct{})

	go func() {
		exp := make(map[string]Dependency)
		exp["circ-b"] = Dependency{Name: "circ-b", Dependencies: []string{"circ-a", "non-yay", "minecraft"}}
		exp["non-yay"] = Dependency{Name: "non-yay", Dependencies: []string{"minecraft"}}
		exp["minecraft"] = Dependency{Name: "minecraft", Dependencies: []string{}}

		actual := make([]Dependency, 0, len(exp))
		r := NewResolverWithFunction("", testInfoResolver)

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

func TestWithCache(t *testing.T) {
	t.Skip("Only used for manual testing at this point in time")

	r := NewResolver()
	it, err := r.Resolve("lib32-eudev")

	assert.NoError(t, err)

	for it.Next() {
		fmt.Printf("Found dependency: %v\n", it.Item())
	}
}

func verify(exp map[string]Dependency, actual []Dependency, t *testing.T) {
	assert.Equal(t, len(exp), len(actual))

	for _, dep := range actual {
		eDep, ok := exp[dep.Name]
		assert.True(t, ok)

		sort.Strings(eDep.Dependencies)
		sort.Strings(dep.Dependencies)
		assert.EqualValues(t, eDep.Dependencies, dep.Dependencies)
	}
}
