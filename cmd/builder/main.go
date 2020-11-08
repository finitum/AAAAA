package main

import (
	"encoding/json"
	"github.com/finitum/AAAAA/pkg/executor"
	"github.com/finitum/AAAAA/pkg/git"
	"github.com/finitum/AAAAA/pkg/makepkg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	cfgStr := os.Getenv("CONFIG")
	if cfgStr == "" {
		log.Fatal("No CONFIG env provided")
	}

	var cfg executor.Config
	err := json.Unmarshal([]byte(cfgStr), &cfg)
	if err != nil {
		log.Fatal("CONFIG env malformed")
	}

	dir, err := ioutil.TempDir(os.TempDir(), "AAAAA_Builder")
	if err != nil {
		log.Fatalf("Couldn't create temp dir %s", dir)
	}

	log.Println("Cloning repo...")
	// Git clone pkg.RepoURL --depth=1
	repo, err := git.Clone(dir, cfg.Package.RepoURL, cfg.Package.RepoBranch)
	if err != nil {
		log.Fatalf("Couldn't clone PGKBUILD repo %v", err)
	}

	hash, err := repo.Head()
	if err != nil {
		log.Fatalf("Getting HEAD failed: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		log.Fatal("Couldn't cd into git repo")
	}

	log.Println("Building pacakge")
	if err := makepkg.Build(); err != nil {
		log.Fatalf("Building package failed %v", err)
	}

	// ls *.pkg.*
	files, err := ioutil.ReadDir(".")
	var found string
	for _, file := range files {
		fname := file.Name()

		if strings.Contains(fname, ".pkg.") && !strings.HasSuffix(fname, ".sig") {
			found = file.Name()
			break
		}
	}

	file, err := os.Open("./.SRCINFO")
	if err != nil {
		log.Fatalf("Couldn't open source info %v", err)
	}

	srcinfo, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Couldn't read source info %v", err)
	}

	src := makepkg.ParseSrcInfo(string(srcinfo))

	log.Printf("parsed srcinfo: %v", src)

	// Upload built package
	log.Println("Uploading package")
	UploadPackage(cfg, src, found, hash.String())
}

func UploadPackage(cfg executor.Config, srcinfo *makepkg.SrcInfo, filename, hash string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Couldn't open package file %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, cfg.UploadURL+"/"+srcinfo.PackageName + "/upload", file)
	if err != nil {
		log.Fatal("yikes2")
	}

	q := req.URL.Query()

	q.Set("waaaa", "luigi")
	q.Set("hash", hash)
	q.Set("filename", filename)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Content-Type", "application/octet-stream")

	log.Printf("Upload URL: %s\n", req.URL.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		log.Fatalf("Couldn't upload file (yikes) %v", err)
	}
}
