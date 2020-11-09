package dependency

import (
	"fmt"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/pkg/errors"
)

type Dependency struct {
	name         string
	dependencies []string
}

type Resolver interface {
	Resolve(pkg string) (*Iterator, error)
}

type aurResolver struct {
	deps         map[string]Dependency
	infoResolver aur.InfoResolveFunction
}

func NewResolver() Resolver {
	return NewResolverWithFunction(aur.SendInfoRequest)
}

func NewResolverWithFunction(infoResolver aur.InfoResolveFunction) Resolver {
	return &aurResolver{
		deps:         make(map[string]Dependency),
		infoResolver: infoResolver,
	}
}

func (r *aurResolver) resolveInternal(pkg string) error {
	res, err := r.infoResolver(pkg)
	if err != nil {
		return errors.Wrap(err, "unable to query aur for dependencies")
	}

	if res.ResultCount > 1 {
		return errors.New(fmt.Sprintf("too many results from aur, expected 1 got %v", res.ResultCount))
	}

	// Not found, so not a dep
	if res.ResultCount == 0 {
		return nil
	}

	deps := append(res.Results[0].Depends, res.Results[0].MakeDepends...)

	// Create dependency record
	r.deps[pkg] = Dependency{
		name:         pkg,
		dependencies: deps,
	}

	// Go through all dependencies
	for _, dep := range deps {
		if _, exists := r.deps[dep]; !exists {

			if err := r.resolveInternal(dep); err != nil {
				return errors.Wrapf(err, "failed to resolve dependency %s", dep)
			}
		}
	}

	return nil
}

// Simple O(n * x) -where x is a factor of hash function efficiency- function to get the intersection of the given
// dependencies with the known AUR dependencies.
func (r *aurResolver) filterDependencies(deps []string) []string {
	res := make([]string, 0, len(r.deps))

	// Do intersection
	for _, dep := range deps {
		if _, ok := r.deps[dep]; ok {
			res = append(res, dep)
		}
	}

	return res
}

func (r *aurResolver) Resolve(pkg string) (*Iterator, error) {
	if err := r.resolveInternal(pkg); err != nil {
		return nil, errors.Wrapf(err, "failed to resolve dependencies for %s", pkg)
	}

	res := make([]Dependency, 0, len(r.deps))
	for _, dep := range r.deps {
		if dep.name != pkg {
			dep.dependencies = r.filterDependencies(dep.dependencies)
			res = append(res, dep)
		}
	}

	it, err := NewIterator(res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create iterator")
	}

	return it, nil
}
