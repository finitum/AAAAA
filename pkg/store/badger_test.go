package store

import (
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestBadger_AddGetDelPackage(t *testing.T) {
	tstPkg := models.Pkg{
		Name:            "name",
		RepoURL:         "url",
		RepoBranch:      "branch",
		KeepLastN:       42,
		LastHash:        plumbing.Hash{'5'},
		UpdateFrequency: time.Hour,
	}

	storePath := os.TempDir() + "/TestBadger_AddPackage"
	store, err := OpenBadgerStore(storePath)
	assert.NoError(t, err)

	err = store.AddPackage(&tstPkg)
	assert.NoError(t, err)

	res, err := store.GetPackage(tstPkg.Name)
	assert.NoError(t, err)

	assert.Equal(t, &tstPkg, res)

	assert.NoError(t, store.DelPackage(&tstPkg))

	assert.NoError(t, os.RemoveAll(storePath))
}

func TestBadger_AllPackages(t *testing.T) {
	tstPkg := models.Pkg{
		Name:            "foo",
		RepoURL:         "url",
		RepoBranch:      "branch",
		KeepLastN:       42,
		LastHash:        plumbing.Hash{'5'},
		UpdateFrequency: time.Hour,
	}

	tstPkg2 := models.Pkg{
		Name:            "bar",
		RepoURL:         "url",
		RepoBranch:      "branch",
		KeepLastN:       42,
		LastHash:        plumbing.Hash{'5'},
		UpdateFrequency: time.Hour,
	}

	storePath := os.TempDir() + "/TestBadger_AllPackages"
	store, err := OpenBadgerStore(storePath)
	assert.NoError(t, err)

	assert.NoError(t, store.AddPackage(&tstPkg))
	assert.NoError(t, store.AddPackage(&tstPkg2))

	pkgs, err := store.AllPackages()
	assert.NoError(t, err)
	assert.Len(t, pkgs, 2)

	pkgNames, err := store.AllPackageNames()
	assert.NoError(t, err)
	assert.Len(t, pkgNames, 2)

	assert.NoError(t, os.RemoveAll(storePath))

}

func TestBadger_AddGetDelUser(t *testing.T) {
	tstUser := models.User{
		Username: "testkees",
		Password: "encryptedyes?",
	}

	storePath := os.TempDir() + "/TestBadger_AddUser"
	store, err := OpenBadgerStore(storePath)
	assert.NoError(t, err)

	err = store.AddUser(&tstUser)
	assert.NoError(t, err)

	res, err := store.GetUser(tstUser.Username)
	assert.NoError(t, err)

	assert.Equal(t, &tstUser, res)

	assert.NoError(t, store.DelUser(&tstUser))

	assert.NoError(t, os.RemoveAll(storePath))
}
