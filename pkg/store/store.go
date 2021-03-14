package store

import (
	"time"

	"github.com/pkg/errors"

	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/finitum/AAAAA/pkg/models"
)

// ErrNotExists is the error returned by the Store or Cache if an entry can not be found, this is useful because it
// isn't always considered an error if no entry exists.
var ErrNotExists = errors.New("entry does not exist")

// Store is a combined interface of PackageStore and UserStore, this is the store the control_plane needs.
type Store interface {
	PackageStore
	UserStore
}

const pkgPrefix = "pkg_"

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

const userPrefix = "user_"

type UserStore interface {
	// GetUser gets a user from the store MUST return ErrNotExists if the user does not exist
	GetUser(name string) (*models.User, error)
	// AddUser adds a user to the store
	AddUser(user *models.User) error
	// DelUser removes a user from the store
	DelUser(username string) error
	// AllUsers lists  all users in the store
	AllUsers() ([]*models.User, error)
	// AllUserNames lists the usernames of all users in the store
	AllUserNames() ([]string, error)
}

// cacheTTL is the TTL of cache entries.
const (
	cacheTTL      = 30 * time.Minute
	resultsPrefix = "results_"
	infoPrefix    = "info_"
)

// Cache interface for caching aur rpc results.
type Cache interface {
	// SetResultsEntry, add an aur result cache entry
	SetResultsEntry(term string, result aur.Results) error
	// GetResultsEntry, retrieve a previously inserted entry (not guaranteed to work due to eviction)
	// returns ErrNotExists if term can't be found
	GetResultsEntry(term string) (aur.Results, error)

	// SetInfoEntry, add an aur info cache entry
	SetInfoEntry(name string, entry *aur.InfoResult) error

	// GetInfoEntry, retrieve a previously inserted entry (not guaranteed to work due to eviction)
	// returns ErrNotExists if term can't be found
	GetInfoEntry(name string) (*aur.InfoResult, error)
}

// GetPartialCacheEntry gets a cache entry even if the term only partially matches the prefix.
func GetPartialCacheEntry(cache Cache, term string) (aur.Results, bool, error) {
	exact := true

	// Keep cutting of letters at the end
	for i := len(term); i > 2; i-- {
		toSearch := term[:i]

		item, err := cache.GetResultsEntry(toSearch)
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
