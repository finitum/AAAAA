/*
Package dependency can be used to resolve Dependencies of packages. The default resolver returned by NewResolver()
will return a resolver that connects to a local aur cache at http://localhost:5001 to resolve Dependencies.

Basic example:

	package main

	import (
		"fmt"
		"github.com/finitum/AAAAA/pkg/dependency"
	)

	func main() {
		r := dependency.NewResolver()

		it, err := r.Resolve("lib32-eudev")
		if err != nil {
			panic(err)
		}

		for it.Next() {
			fmt.Printf("Found dependency: %v\n", it.Item())
		}
	}
*/
package dependency

import (
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/pkg/errors"
)

// Dependency represents a dependency of an AUR package.
// This dependency has a name, and a slice of name of dependencies of this dependency.
type Dependency struct {
	Name         string
	Dependencies []string
}

/*
Resolver represents a resolver, which is able to resolve all dependencies of a given package with the given name.

Most use cases can be solved by using a custom InfoResolveFunction in combination with a custom URL, but it is possible
to provide an alternative implementation of the default resolver by implementing this interface.
*/
type Resolver interface {
	// Resolve accepts a string representing the package name, and will return an Iterator, which returns all
	// (in)direct dependencies of the given name in no particular order. Furthermore, it is legal for the Iterator
	// to contain a set of dependencies which have a circular relation, the consumer should handle this case.
	Resolve(pkg string) (*Iterator, error)
}

type aurResolver struct {
	url          string
	deps         map[string]Dependency
	infoResolver aur.InfoResolveFunction
}

// NewResolver returns a default Resolver. This resolver uses aur.SendCachedInfoRequest to make requests to the
// default url 'http://localhost:5001/info/%s'.
func NewResolver() Resolver {
	return NewResolverWithFunction("http://localhost:5001/info/%s", aur.SendCachedInfoRequest)
}

// NewResolverWithFunction returns a Resolver, which uses the given url. The url format depends on the given
// aur.InfoResolveFunction, but for the default function aur.SendCachedInfoRequest the url should contain one '%s',
// which is filled with the current package name.
func NewResolverWithFunction(url string, infoResolver aur.InfoResolveFunction) Resolver {
	return &aurResolver{
		url:          url,
		deps:         make(map[string]Dependency),
		infoResolver: infoResolver,
	}
}

func (r *aurResolver) resolveInternal(pkg string) error {
	res, err := r.infoResolver(r.url, pkg)

	// Not found so not a dependency
	if err == aur.NotInAurErr {
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "unable to query aur for dependencies")
	}

	deps := append(res.Depends, res.MakeDepends...)

	// Create dependency record
	r.deps[pkg] = Dependency{
		Name:         pkg,
		Dependencies: deps,
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
		if dep.Name != pkg {
			dep.Dependencies = r.filterDependencies(dep.Dependencies)
			res = append(res, dep)
		}
	}

	it, err := NewIterator(res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create iterator")
	}

	return it, nil
}
