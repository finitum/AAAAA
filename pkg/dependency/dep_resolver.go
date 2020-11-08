package dependency

import (
	"fmt"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/pkg/errors"
)

type Resolver interface {
	Resolve(pkg string) (*Iterator, error)
}

type aurResolver struct {
	deps         map[string]struct{}
	infoResolver aur.InfoResolveFunction
}

func NewResolver() Resolver {
	return NewResolverWithFunction(aur.SendInfoRequest)
}

func NewResolverWithFunction(infoResolver aur.InfoResolveFunction) Resolver {
	return &aurResolver{
		deps:         make(map[string]struct{}),
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

	r.deps[pkg] = struct{}{}

	// Go through all dependencies
	for _, dep := range append(res.Results[0].Depends, res.Results[0].MakeDepends...) {
		if _, exists := r.deps[dep]; !exists {

			if err := r.resolveInternal(dep); err != nil {
				return errors.Wrapf(err, "failed to resolve dependency %s", dep)
			}
		}
	}

	return nil
}

func (r *aurResolver) Resolve(pkg string) (*Iterator, error) {
	if err := r.resolveInternal(pkg); err != nil {
		return nil, errors.Wrapf(err, "failed to resolve dependencies for %s", pkg)
	}

	// Delete itself from dependency list
	delete(r.deps, pkg)

	// TODO(timanema): Return iterator instead of []string
	res := make([]string, 0, len(r.deps))
	for dep := range r.deps {
		res = append(res, dep)
	}

	it, err := NewIterator(res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create iterator")
	}

	return it, nil
}
