package store

import (
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/pkg/errors"
	"time"
)

// ErrNotExists is the error returned by the Store or Cache if an entry can not be found, this is useful because it
// isn't always considered an error if no entry exists.
var ErrNotExists = errors.New("entry does not exist")

// Store is a combined interface of PackageStore and UserStore, this is the store the control_plane needs
type Store interface {
	PackageStore
	UserStore
}

type PackageStore interface {
	// GetPackage gets a package definition from the store MUST return ErrNotExists if the package does not exist
	GetPackage(name string) (*models.Pkg, error)
	// AddPackage add a package to the store
	AddPackage(pkg *models.Pkg) error
	// DelPackage remove a package from the store
	DelPackage(name string) error
	// AllPackages lists all packages inside the store
	AllPackages() ([]*models.Pkg, error)
	// AllPackageNames lists the names of all the packages in the store
	AllPackageNames() ([]string, error)
}

type UserStore interface {
	// GetUser gets a user from the store MUST return ErrNotExists if the user does not exist
	GetUser(name string) (*models.User, error)
	// AddUser adds a user to the store
	AddUser(user *models.User) error
	// DelUser removes a user from the store
	DelUser(user *models.User) error
	// AllUserNames lists the usernames of all users in the store
	AllUserNames() ([]string, error)
}

// cacheTTL is the TTL of cache entries
const cacheTTL = 30 * time.Minute

// Cache interface for caching aur rpc results
type Cache interface {
	// SetEntry, add an aur result cache entry
	SetEntry(term string, result aur.Results) error
	// GetEntry, retrieve a previously inserted entry (not guaranteed to work due to eviction)
	// returns ErrNotExists if term can't be found
	GetEntry(term string) (aur.Results, error)
}

// GetPartialCacheEntry gets a cache entry even if the term only partially matches the prefix
func GetPartialCacheEntry(cache Cache, term string) (aur.Results, bool, error) {
	exact := true

	// Keep cutting of letters at the end
	for i := len(term); i > 2; i-- {
		toSearch := term[:i]

		item, err := cache.GetEntry(toSearch)
		if err == ErrNotExists {
			// continue searching if it doesn't match
			exact = false
			continue
		} else if err != nil {
			return nil, false, errors.Wrap(err, "getting partial entry")
		}

		return item, exact, nil
	}

	return nil, false, ErrNotExists
}
