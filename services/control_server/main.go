package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
)

func main() {
	// Open Database
	db, err := store.OpenBadgerStore(os.TempDir() + "/AAAAA")
	if err != nil {
		log.Fatal(err)
	}

	tokenAuth := jwtauth.New(jwt.SigningMethodHS384.Name, []byte("change me"), nil)

	// Router
	rs := NewRoutes(db, tokenAuth)

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Compress(5))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", rs.HelloWorld)

	// Protected Routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens.
		r.Use(jwtauth.Authenticator)

	})

	log.Fatal(http.ListenAndServe(":5000", r))
}
