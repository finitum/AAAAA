package store

import "github.com/finitum/AAAAA/pkg/models"

type PackageStore interface {
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

type Store interface {
	PackageStore
	UserStore
}
