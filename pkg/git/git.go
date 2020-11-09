package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
)

func LatestHash(url, branch string) (plumbing.Hash, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{url},
		Fetch: nil,
	})

	rfs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return plumbing.Hash{}, errors.Wrap(err, "listing remote")
	} else if len(rfs) < 1 {
		return plumbing.Hash{}, errors.New("no references on git repo")
	}

	refName := plumbing.NewBranchReferenceName(branch)

	var foundRef *plumbing.Reference
	for _, ref := range rfs {
		if ref.Name() == refName {
			foundRef = ref
			break
		}
	}
	if foundRef == nil {
		return [20]byte{}, errors.New("no ref found")
	}

	return foundRef.Hash(), nil
}

func Clone(path, url, branch string) (*git.Repository, error) {
	refName := plumbing.NewBranchReferenceName(branch)

	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               url,
		RemoteName:        "origin",
		ReferenceName:     refName,
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: 1,
	})
	if err != nil {
		return nil, errors.Wrap(err, "git clone")
	}

	return repo, nil
}

func InMemoryClone(url, branch string) (*git.Repository, error) {
	refName := plumbing.NewBranchReferenceName(branch)

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:               url,
		RemoteName:        "origin",
		ReferenceName:     refName,
		SingleBranch:      true,
		Depth:             1,
		RecurseSubmodules: git.NoRecurseSubmodules,
	})
	if err != nil {
		return nil, errors.Wrap(err, "git clone")
	}

	return repo, nil
}
