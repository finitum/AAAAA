package repo_add

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestNewRepoAdd(t *testing.T) {
	_, err := NewRepoAdd("test")

	// If NewRepoAdd fails, the tests don't run on an arch system.
	// By returning the test just passes on non Archlinux systems
	if err != nil {
		t.Skip("Not running on an archlinux system, skipping.")
	}
}

func TestIntegration(t *testing.T) {
	dir := os.TempDir() + "/AAAAA_TestIntegration"
	dbpath := dir + "/test.db.tar.gz"
	pkgbuildpath := dir + "/PKGBUILD"

	repo, err := NewRepoAdd(dbpath)

	// If NewRepoAdd fails, the tests don't run on an arch system.
	// By returning the test just passes on non Archlinux systems
	if err != nil {
		t.Skip("Not running on an archlinux system, skipping.")
	}

	// create tempdir to run in
	err = os.Mkdir(dir, os.ModePerm)
	assert.NoError(t, err)

	// Create a dummy packagebuild
	file, err := os.Create(pkgbuildpath)
	assert.NoError(t, err)

	_, err = file.WriteString(`
pkgname='dummy'
pkgver=0.1
pkgrel=1
arch=(any)
provides=('dummy=0.1')`)
	assert.NoError(t, err)

	// Build this dummy package
	err = os.Chdir(dir)
	assert.NoError(t, err)
	cmd := exec.Command("makepkg", "-cf", "PKGDEST="+dir, "PKGEXT=.pkg.tar.gz")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	assert.NoError(t, err)

	err = repo.AddPackage(dir+"/dummy-0.1-1-any.pkg.tar.gz", RepoAddOptions{})
	assert.NoError(t, err)

	// Assert that it created the database file
	f, err := os.Stat(dbpath)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	// TODO: This test should actually assert things about the database
	// 		 I did manually check it working, but that't not really enough is it :P

	err = os.RemoveAll(dir)
	assert.NoError(t, err)
}
