package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
)

// Pkg is an archlinux package.
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
	LastHash plumbing.Hash `json:",omitempty"`
	// UpdateFrequency determines how often the package should be updated
	UpdateFrequency time.Duration

	// TODO: Version?
}

// Render will run before marshalling a Pkg, good place to do pre-processing.
func (p *Pkg) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

// Bind will run after unmarshalling a Pkg, good place to do post-processing.
func (p *Pkg) Bind(*http.Request) error {
	if p.Name == "" || p.RepoURL == "" {
		return errors.New("package is missing required fields Name and/or RepoURL")
	}

	if p.RepoBranch == "" {
		p.RepoBranch = "master"
	}

	if p.KeepLastN == 0 {
		p.KeepLastN = 2
	}

	if p.UpdateFrequency == 0 {
		p.UpdateFrequency = time.Hour
	}

	return nil
}

type User struct {
	Username string
	Password string
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	u.Password = ""

	return nil
}

func (u *User) Bind(*http.Request) error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid user")
	}

	return nil
}
