package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/finitum/AAAAA/pkg/executor"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/finitum/AAAAA/services/control_server/config"
	"github.com/finitum/AAAAA/services/control_server/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg := config.Default()
	if err := cfg.CreateDirectories(); err != nil {
		log.Fatalf("Couldn't create directories %v", err)
	}

	// Open Database
	db, err := store.OpenBadger(cfg.StoreLocation)
	if err != nil {
		log.Fatalf("Opening Badger store failed: %v", err)
	}
	defer db.Close()

	tokenAuth := jwtauth.New(jwt.SigningMethodHS384.Name, []byte(cfg.JWTKey), nil)

	// Auth service
	auths := auth.NewStoreAuth(db, tokenAuth)

	// Create initial user
	initialUser(db, auths)

	// Executor
	exec, err := executor.NewDockerExecutor(cfg.RunnerImage)
	if err != nil {
		log.Fatalf("Starting docker executor failed: %v", err)
	}

	// Router
	rs := routes.New(cfg, db, auths, exec)

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.Method == http.MethodOptions {
				return
			}

			h.ServeHTTP(w, r)
		})
	})
	//r.Use(middleware.Compress(5))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", rs.HelloWorld)

	r.Post("/login", rs.Login)
	r.Get("/package", rs.GetPackages)

	// Protected Routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens.
		r.Use(jwtauth.Authenticator)
		//r.Use(corsHandler)

		r.Post("/user", rs.AddUser)
		r.Post("/package", rs.AddPackage)

		r.Delete("/package/{pkg}", rs.RemovePackage)
		r.Put("/package/{pkg}", rs.UpdatePackage)

		r.Post("/package/{pkg}/upload", rs.UploadPackage)
		r.Put("/package/{pkg}/build", rs.TriggerBuild)
	})

	log.Fatal(http.ListenAndServe(cfg.Address, r))
}

func initialUser(db store.Store, auths auth.AuthenticationService) {
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

	if err := auths.Register(&models.User{
		Username: "admin",
		Password: pass,
	}); err != nil {
		log.Fatal(err)
	}

	log.Infof("|> username: admin, password: %s \n", pass)
}
