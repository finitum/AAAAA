package executor

import (
	"context"
	"github.com/finitum/AAAAA/pkg/models"
)

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
	// PrepareBuild should prepare the executor for incoming executions, this can be used to pull the latest
	// docker image for example
	PrepareBuild(ctx context.Context) error

	// BuildPackage should build the latest package and upload it to the url specified in the Config
	BuildPackage(ctx context.Context, cfg *Config) error
}
