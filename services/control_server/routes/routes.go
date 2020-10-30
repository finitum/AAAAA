package routes

import (
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/finitum/AAAAA/pkg/git"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/render"
	"net/http"
)

type Routes struct {
	db   store.Store
	auth auth.AuthenticationService
}

func New(db store.Store, auth auth.AuthenticationService) *Routes {
	return &Routes{db, auth}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

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

// POST - /package
// body: models.Pkg
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
