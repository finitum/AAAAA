package models

import (
	"github.com/go-git/go-git/v5/plumbing"
	"time"
)

// Pkg is an archlinux package
type Pkg struct {
	// Name is the name of the package (unique)
	Name string
	// RepoURL is the git repository where the PKGBUILD can be found
	RepoURL string
	// RepoBranch is which branch is used for updates
	RepoBranch string
	// KeepLastN determines how many old versions of packages are kept
	KeepLastN int
	// LastHash is the latest SHA1 retrieved from the package repo
	LastHash plumbing.Hash
	// UpdateFrequency determines how often the package should be updated
	UpdateFrequency time.Duration

	// TODO: Version?
}

func DefaultPkg() Pkg {
	return Pkg{
		KeepLastN:       2,
		UpdateFrequency: time.Hour,
	}
}

type User struct {
	Username string
	Password string
}
