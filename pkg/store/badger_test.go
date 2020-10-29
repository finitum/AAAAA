package store

import (
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestBadger_AddGetPackage(t *testing.T) {
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

	assert.NoError(t, os.RemoveAll(storePath))
}
