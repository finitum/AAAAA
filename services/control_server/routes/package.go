package routes

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"

	"github.com/finitum/AAAAA/pkg/executor"
	"github.com/finitum/AAAAA/pkg/git"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/repo_add"
	"github.com/finitum/AAAAA/pkg/store"
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
	if err != store.ErrNotExists {
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

func (rs *Routes) RemovePackage(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")

	_, err := rs.db.GetPackage(pkgName)
	if err == store.ErrNotExists {
		_ = render.Render(w, r, ErrNotFound())
		return
	}

	if err := rs.db.DelPackage(pkgName); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}
}

func (rs *Routes) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")

	pkg, err := rs.db.GetPackage(pkgName)
	if err == store.ErrNotExists {
		_ = render.Render(w, r, ErrNotFound())
		return
	}

	if err := render.Bind(r, pkg); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := rs.db.AddPackage(pkg); err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}
}

func (rs *Routes) TriggerBuild(w http.ResponseWriter, r *http.Request) {
	pkgName := chi.URLParam(r, "pkg")

	pkg, err := rs.db.GetPackage(pkgName)
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	tokenStr := token.Raw

	go func() {
		ctx := context.Background()

		if err := rs.exec.PrepareBuild(ctx); err != nil {
			log.Warnf("trigger prepare build %v", err)
		}

		if err := rs.exec.BuildPackage(ctx, &executor.Config{
			Package:   pkg,
			Token:     tokenStr,
			UploadURL: rs.cfg.ExternalAddress + "/package",
		}); err != nil {
			log.Warnf("trigger build %v", err)
			return
		}
	}()

	w.WriteHeader(http.StatusAccepted)
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
	pkgPath := rs.cfg.RepoLocation + "/" + filename
	if externalUrl == "" {
		file, createErr := os.Create(pkgPath)
		if createErr != nil {
			_ = render.Render(w, r, ErrServerError(createErr))
			log.Warnf("UploadPackage creating file failed: %v", createErr)
			return
		}

		// or fs.CopyStream if this returns ErrUnexpectedEof
		if _, copyErr := io.Copy(file, r.Body); copyErr != nil {
			_ = render.Render(w, r, ErrServerError(copyErr))
			log.Warnf("UploadPackage writing to file failed: %v", copyErr)
			return
		}
	} else {
		panic("not implemented!")
	}

	pkg.LastHash = plumbing.NewHash(hash)
	if addPackageErr := rs.db.AddPackage(pkg); addPackageErr != nil {
		_ = render.Render(w, r, ErrServerError(addPackageErr))
		log.Warnf("UploadPackage updating db failed: %v", addPackageErr)
		return
	}

	// TODO: fix
	ra, err := repo_add.NewRepoAdd(rs.cfg.RepoLocation)
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
