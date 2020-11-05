package store

import (
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/pkg/errors"
	"time"
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

const cacheTTL = 30 * time.Minute

type Cache interface {
	SetEntry(searchterm string, result aur.Results) error
	GetEntry(searchterm string) (aur.Results, error)
}

type Store interface {
	PackageStore
	UserStore
}

func GetPartialCacheEntry(cache Cache, term string) (aur.Results, bool, error) {
	exact := true
	for i := len(term); i > 2; i-- {
		toSearch := term[:i]

		item, err := cache.GetEntry(toSearch)
		if err == ErrNotExists {
			exact = false
			continue
		} else if err != nil {
			return nil, false, errors.Wrap(err, "getting partial entry")
		}

		return item, exact, nil
	}

	return nil, false, ErrNotExists
}
