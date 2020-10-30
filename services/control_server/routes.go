package main

import (
	"github.com/finitum/AAAAA/pkg/git"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/render"
	"net/http"
)

type Routes struct {
	db store.Store
}

func NewRoutes(db store.Store) *Routes {
	return &Routes{db}
}

type AppCode int64
const (
	AppCodeGeneric AppCode = iota
	AppCodeGitRepoUnreachable
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string  `json:"status"`          // user-level status message
	AppCode    AppCode `json:"code,omitempty"`  // application-specific error code
	ErrorText  string  `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error, code ...AppCode) render.Renderer {
	retcode := AppCodeGeneric
	if len(code) > 0 {
		retcode = code[0]
	}

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid Request",
		AppCode: 		retcode,
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Server Error",
		ErrorText:      err.Error(),
	}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

// POST - /package
// body: models.Pkg
func (rs *Routes) AddPackage(w http.ResponseWriter, r *http.Request) {
	var pkg models.Pkg

	if err := render.Bind(r, &pkg); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// TODO: existence check?

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
