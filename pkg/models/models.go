package models

import (
	"errors"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Pkg is an archlinux package
type Pkg struct {
	// Name is the name of the package (unique)
	Name string
	// RepoURL is the git repository where the PKGBUILD can be found
	RepoURL string
	// RepoBranch is which branch is used for updates
	RepoBranch string
	// KeepLastN determines how many old versions of packages are kept
	KeepLastN int
	// LastHash is the latest SHA1 retrieved from the package repo
	LastHash plumbing.Hash `json:",omitempty"`
	// UpdateFrequency determines how often the package should be updated
	UpdateFrequency time.Duration

	// TODO: Version?
}

// Render will run before marshalling a Pkg, good place to do pre-processing
func (p *Pkg) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

// Bind will run after unmarshalling a Pkg, good place to do post-processing
func (p *Pkg) Bind(*http.Request) error {
	if p.Name == "" || p.RepoURL == "" {
		return errors.New("package is missing required fields Name and/or RepoURL")
	}

	if p.RepoBranch == "" {
		p.RepoBranch = "master"
	}

	if p.KeepLastN == 0 {
		p.KeepLastN = 2
	}

	if p.UpdateFrequency == 0 {
		p.UpdateFrequency = time.Hour
	}

	return nil
}

type User struct {
	Username string
	Password string
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	u.Password = ""

	return nil
}

func (u *User) Bind(*http.Request) error {
	if u.Username == "" || u.Password == "" {
		return errors.New("invalid user")
	}

	return nil
}

type Job struct {
	PackageName string
	Status      BuildStatus
	Logs        BuildLog `json:",omitempty"`
	Uuid        string
	Time        time.Time
}

// logsToKeep are the number of log lines to keep when sending a job.
const logsToKeep = 10

func (j Job) Render(w http.ResponseWriter, r *http.Request) error {
	// Remove everything but the last 10 log lines. To get all
	// logs the /job/{uuid}/logs route can be used. This is because
	// the logs can get quite large, and if you want information about a single
	// job it's not really useful to get all the logs. This is especially true
	// when retrieving *all* jobs. In that case you really don't want all logs to
	// be sent over as well
	if len(j.Logs) > logsToKeep {
		j.Logs = j.Logs[len(j.Logs)-logsToKeep:]
	}

	return nil
}

type BuildStatus int

const (
	BuildStatusPending BuildStatus = iota
	BuildStatusPullingRepo
	BuildStatusRunning
	BuildStatusUploading
	BuildStatusDone

	BuildStatusErrored
)

type LogLine struct {
	Time    time.Time
	Level   log.Level
	message string
}

func (j LogLine) Bind(r *http.Request) error {
	return nil
}

type BuildLog []LogLine

func (j LogLine) Render(w http.ResponseWriter, r *http.Request) error {

	return nil
}
