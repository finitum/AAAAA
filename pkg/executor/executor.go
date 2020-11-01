package executor

import "github.com/finitum/AAAAA/pkg/models"

type Config struct {
	// Package
	Package *models.Pkg
	// UploadURL is where the executor should POST the finished package to
	UploadURL string
	// Token
	Token string
	// RepoURL
	RepoURL string
}

type Executor interface {
	// BuildPackage should build the latest package and upload it to the url specified in the Config
	BuildPackage(cfg *Config) error
}
