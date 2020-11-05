package config

import "os"

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
}

// Default returns a default set of config values which generally work for local development
func Default() *Config {
	root := "./AAAAA/"
	return &Config{
		StoreLocation:   root + "store",
		RepoLocation:    root + "repo",
		RunnerImage:     "aaaaa-builder",
		JWTKey:          "change-me",
		Address:         "0.0.0.0:5000",
		ExternalAddress: "0.0.0.0:5000",
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
