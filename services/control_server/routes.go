package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/finitum/AAAAA/pkg/git"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Routes struct {
	db        store.Store
	tokenAuth *jwtauth.JWTAuth
}

func NewRoutes(db store.Store, auth *jwtauth.JWTAuth) *Routes {
	return &Routes{db, auth}
}

func (*Routes) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello World!"))
}

func (rs *Routes) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := render.Bind(r, &user); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	dbUser, err := rs.db.GetUser(user.Username)
	if err != nil {
		_ = render.Render(w, r, ErrUnauthorized())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		_ = render.Render(w, r, ErrUnauthorized())
		return
	}

	_, tokenString, err := rs.tokenAuth.Encode(jwt.MapClaims{"username": dbUser.Username})
	if err != nil {
		_ = render.Render(w, r, ErrServerError(err))
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
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
