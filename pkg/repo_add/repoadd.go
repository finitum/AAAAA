// Package repo_add provides a wrapper around the
// repo_add program.
package repo_add

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const repoAddCommand string = "repo-add"
const repoRemoveCommand string = "repo-remove"

// CommonOptions are options that are common to the repo-add and repo-remove commands.
// Derived from https://www.archlinux.org/pacman/repo-add.8.html
type CommonOptions struct {
	// Generate a PGP signature file using GnuPG.
	// This will execute gpg --detach-sign --use-agent on the generated database to generate
	// a detached signature file, using the GPG agent if it is available.
	// The signature file will be the entire filename of the database with a “.sig” extension.
	// equivalent to `repo-add --sign`
	Sign bool

	// Specify a key to use when signing packages.
	// Can also be specified using the GPGKEY environmental variable.
	// If not specified in either location, the default key from the keyring will be used.
	// equivalent to `repo-add --key`
	Key string

	// Verify the PGP signature of the database before updating the database.
	// If the signature is invalid, an error is produced and the update does not proceed.
	// equivalent to `repo-add --verify`
	Verify bool
}

// RepoAddOptions are options which are only used for the repo-add command.
// Includes the common options.
// Derived from https://www.archlinux.org/pacman/repo-add.8.html
type RepoAddOptions struct {
	CommonOptions

	// Remove old package files from the disk when updating their entry in the database.
	// equivalent to `repo-add --remove`
	RemoveOld bool

	// Only add packages that are not already in the database.
	// Warnings will be printed upon detection of existing packages, but they will not be re-added.
	// equivalent to `repo-add --new`
	OnlyNew bool
}

// RepoAdd has two functions which do the same as the `repo-add` and `repo-remove` commands
// on Archlinux. The interface of the functions is derived from https://www.archlinux.org/pacman/repo-add.8.html
type RepoAdd struct {
	dbpath string
}

// NewRepoAdd makes a new RepoAdd struct.
// It also makes sure the repo-add script exists on the system.
func NewRepoAdd(dbpath string) (*RepoAdd, error) {
	_, err := exec.LookPath(repoAddCommand)
	if err != nil {
		return nil, errors.New("the repo-add binary must be present to run AAAAA")
	}

	return &RepoAdd{dbpath}, nil
}

func serializeCommonOptions(options CommonOptions) []string {
	// 8 because that's the maximum length it can ever be in the
	// serializeCommonOptions, serializeAddOptions, Add and Remove
	// functions combined
	res := make([]string, 0, 8)

	if options.Sign {
		res = append(res, "--sign")
	}

	if options.Verify {
		res = append(res, "--verify")
	}

	if options.Key != "" {
		res = append(res, "--key", options.Key)
	}

	return res
}

func serializeAddOptions(options RepoAddOptions) []string {
	res := serializeCommonOptions(options.CommonOptions)

	if options.OnlyNew {
		res = append(res, "--new")
	}

	if options.RemoveOld {
		res = append(res, "--remove")
	}

	return res
}

// AddPackage is a wrapper around the repo-add command.
func (r *RepoAdd) AddPackage(packagepath string, options RepoAddOptions) error {
	so := serializeAddOptions(options)

	so = append(so, r.dbpath)
	so = append(so, packagepath)

	cmd := exec.Command(repoAddCommand, so...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	println(cmd.String())

	err := cmd.Run()
	if err != nil {
		return err
	}

	if !cmd.ProcessState.Success() {
		return fmt.Errorf("running repo-add failed (%d)", cmd.ProcessState.ExitCode())
	}

	return nil
}

// RemovePackage is a wrapper around the repo-add command.
func (r *RepoAdd) RemovePackage(packagepath string, options CommonOptions) error {
	so := serializeCommonOptions(options)

	so = append(so, r.dbpath)
	so = append(so, packagepath)

	cmd := exec.Command(repoAddCommand, so...)
	cmd.Args[0] = repoRemoveCommand

	err := cmd.Run()
	if err != nil {
		return err
	}

	if !cmd.ProcessState.Success() {
		return errors.New("running repo-remove failed")
	}

	return nil
}
