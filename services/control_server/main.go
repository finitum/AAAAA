package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

func main() {
	// Open Database
	db, err := store.OpenBadgerStore(os.TempDir() + "/AAAAA")
	if err != nil {
		log.Fatal(err)
	}

	// Create initial user
	initialUser(db)

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

	r.Post("/login", rs.Login)

	// Protected Routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens.
		r.Use(jwtauth.Authenticator)

	})

	log.Fatal(http.ListenAndServe(":5000", r))
}

func initialUser(db store.Store) {
	users, err := db.AllUserNames()
	if err != nil {
		log.Fatal(err)
	}
	if len(users) != 0 {
		return
	}

	log.Info("Creating default admin user as no users were found")
	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	pass := base64.StdEncoding.EncodeToString(buf)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Username: "admin",
		Password: string(hashedPass),
	}

	if err := db.AddUser(&user); err != nil {
		log.Fatal(err)
	}

	log.Infof(" username: admin, password: %s \n", pass)
}
