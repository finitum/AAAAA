package config

type Config struct {
	StoreLocation string
	RepoLocation  string
	RunnerImage   string
	JWTKey        string
	Address       string
}

func Default() *Config {
	root := "./AAAAA/"
	return &Config{
		StoreLocation: root + "store",
		RepoLocation:  root + "repo",
		RunnerImage:   "aaaaa-builder",
		JWTKey:        "change-me",
		Address:       "0.0.0.0:5000",
	}
}
