package store

import (
	"errors"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/finitum/AAAAA/pkg/models"
)

var ErrNotExists = errors.New("entry does not exist")

type PackageStore interface {
	// GetPackage gets a package definition from the store MUST return ErrNotExists if the package does not exist
	GetPackage(name string) (*models.Pkg, error)
	AddPackage(pkg *models.Pkg) error
	DelPackage(pkg *models.Pkg) error
	AllPackages() ([]*models.Pkg, error)
	AllPackageNames() ([]string, error)
}

type UserStore interface {
	GetUser(name string) (*models.User, error)
	AddUser(user *models.User) error
	DelUser(user *models.User) error
	AllUserNames() ([]string, error)
}

type Cache interface {
	SetEntry(searchterm string, result aur.Results) error
	GetEntry(searchterm string) (aur.Results, bool, error)
}

type Store interface {
	PackageStore
	UserStore
}
