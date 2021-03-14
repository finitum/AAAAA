package store

import (
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"

	"github.com/finitum/AAAAA/pkg/models"
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
	store, err := OpenBadger(storePath)
	assert.NoError(t, err)

	err = store.AddPackage(&tstPkg)
	assert.NoError(t, err)

	res, err := store.GetPackage(tstPkg.Name)
	assert.NoError(t, err)

	assert.Equal(t, &tstPkg, res)

	assert.NoError(t, store.DelPackage(tstPkg.Name))

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
	store, err := OpenBadger(storePath)
	assert.NoError(t, err)

	assert.NoError(t, store.AddPackage(&tstPkg))
	assert.NoError(t, store.AddPackage(&tstPkg2))

	pkgs, err := store.AllPackages()
	assert.NoError(t, err)
	assert.Len(t, pkgs, 2)

	assert.Contains(t, pkgs, &tstPkg, &tstPkg2)

	pkgNames, err := store.AllPackageNames()
	assert.NoError(t, err)
	assert.Len(t, pkgNames, 2)

	assert.Contains(t, pkgNames, tstPkg.Name, tstPkg2.Name)

	assert.NoError(t, os.RemoveAll(storePath))
}

func TestBadger_AddGetDelUser(t *testing.T) {
	tstUser := models.User{
		Username: "testkees",
		Password: "encryptedyes?",
	}

	storePath := os.TempDir() + "/TestBadger_AddUser"
	store, err := OpenBadger(storePath)
	assert.NoError(t, err)

	err = store.AddUser(&tstUser)
	assert.NoError(t, err)

	res, err := store.GetUser(tstUser.Username)
	assert.NoError(t, err)

	assert.Equal(t, &tstUser, res)

	assert.NoError(t, store.DelUser(tstUser.Username))

	assert.NoError(t, os.RemoveAll(storePath))
}

func TestBadger_AllUserNames(t *testing.T) {
	tstUser := models.User{
		Username: "yoink",
		Password: "rot26bestencryption",
	}
	tstUser2 := models.User{
		Username: "; DROP TABLE 'USERS';--",
		Password: "2xrot13bestencryption",
	}

	storePath := os.TempDir() + "/TestBadger_AllUserNames"
	store, err := OpenBadger(storePath)
	assert.NoError(t, err)

	assert.NoError(t, store.AddUser(&tstUser))
	assert.NoError(t, store.AddUser(&tstUser2))

	names, err := store.AllUserNames()
	assert.NoError(t, err)
	assert.Len(t, names, 2)

	assert.Contains(t, names, tstUser.Username, tstUser2.Username)

	assert.NoError(t, os.RemoveAll(storePath))
}
