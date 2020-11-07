package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

// Config is the configuration used by the control server
type Config struct {
	// StoreLocation is the Location where the database is stored
	StoreLocation string
	// RepoLocation is the Location where the pacman repository is stored
	RepoLocation string
	// RunnerImage is the docker image used to build packages
	RunnerImage string
	// JWTKey private key for signing JWTs
	JWTKey string
	// Address to bind the http server on
	Address string
	// ExternalAddress the address which is externally reachable (used for the runner amongst others
	ExternalAddress string
	// Executor determines where packages are built can be either 'docker' or 'kubernetes'
	Executor string
	// KubeConfigPath points to a kube config if using the kubernetes executor in external mode, leave empty for internal mode
	KubeConfigPath string `json:",omitempty"`
	// KubeNamespace is the namespace in which the kubernetes jobs should be ran, defaults to kubernetes' default
	KubeNamespace string `json:",omitempty"`
}

// DevDefault returns a default set of config values which generally work for local development
func DevDefault() *Config {
	root := "./AAAAA/"

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("Couldn't get user's homedir %v", err)
		homedir = root
	}

	return &Config{
		StoreLocation:   root + "store",
		RepoLocation:    root + "repo",
		RunnerImage:     "harbor.xirion.net/library/aaaaa-builder",
		JWTKey:          "change-me",
		Address:         "0.0.0.0:5000",
		ExternalAddress: "0.0.0.0:5000",
		Executor:        "docker",
		KubeConfigPath:  path.Join(homedir, ".kube", "config"),
		KubeNamespace:   "",
	}
}

func (cfg *Config) CreateDirectories() error {
	if err := os.MkdirAll(cfg.StoreLocation, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.RepoLocation, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// TODO: Config from env
