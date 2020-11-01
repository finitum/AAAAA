package routes

import (
	"errors"
	"github.com/finitum/AAAAA/pkg/git"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/repo_add"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func (rs *Routes) GetPackages(w http.ResponseWriter, r *http.Request) {
	dbPkgs, err := rs.db.AllPackages()
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
	}

	pkgs := make([]render.Renderer, len(dbPkgs))
	for i, pkg := range dbPkgs {
		pkgs[i] = pkg
	}

	_ = render.RenderList(w, r, pkgs)
}

func (rs *Routes) AddPackage(w http.ResponseWriter, r *http.Request) {
	var pkg models.Pkg

	if err := render.Bind(r, &pkg); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	_, err := rs.db.GetPackage(pkg.Name)
	if err == store.ErrNotExists {
		_ = render.Render(w, r, ErrExists())
		return
	}

	// Check if we are able to ls-remote the repo
	if _, err := git.LatestHash(pkg.RepoURL, pkg.RepoBranch); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err, AppCodeGitRepoUnreachable))
		return
	}

	if err := rs.db.AddPackage(&pkg); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rs *Routes) UploadPackage(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")
	hash := r.URL.Query().Get("hash")
	externalUrl := r.URL.Query().Get("remote_url")
	filename := r.URL.Query().Get("filename")

	if pkgName == "" || hash == "" || filename == "" {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("invalid query param")))
		return
	}

	// update latest hash?
	pkg, err := rs.db.GetPackage(pkgName)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	if pkg.LastHash.String() == hash {
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("hash is already latest")))
		return
	}

	// Download file
	pkgPath := "./AAAAA/repo/" + filename
	if externalUrl == "" {
		file, err := os.Create(pkgPath)
		if err != nil {
			_ = render.Render(w, r, ErrServerError(err))
			log.Warnf("UploadPackage creating file failed: %v", err)
			return
		}

		// or fs.CopyStream if this returns ErrUnexpectedEof
		if _, err := io.Copy(file, r.Body); err != nil {
			_ = render.Render(w, r, ErrServerError(err))
			log.Warnf("UploadPackage writing to file failed: %v", err)
			return
		}
	} else {
		panic("not implemented!")
	}

	pkg.LastHash = plumbing.NewHash(hash)
	if err := rs.db.AddPackage(pkg); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Warnf("UploadPackage updating db failed: %v", err)
		return
	}

	// repo-add
	ra, err := repo_add.NewRepoAdd("./AAAAA")
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Warnf("UploadPackage repo add failed: %v", err)
		return
	}

	if err := ra.AddPackage(pkgPath, &repo_add.RepoAddOptions{}); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		log.Warnf("UploadPackage repo add failed: %v", err)
		return
	}

	// TODO: sign


	w.WriteHeader(http.StatusCreated)
}
