package main

import (
	"github.com/finitum/AAAAA/internal/cors"
	"github.com/finitum/AAAAA/pkg/auth"
	"github.com/finitum/AAAAA/pkg/executor"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/finitum/AAAAA/services/control_server/config"
	"github.com/finitum/AAAAA/services/control_server/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg := config.DevDefault()
	if err := cfg.CreateDirectories(); err != nil {
		log.Fatalf("Couldn't create directories %v", err)
	}

	// Open Database
	db, err := store.OpenBadger(cfg.StoreLocation)
	if err != nil {
		log.Fatalf("Opening Badger store failed: %v", err)
	}
	defer db.Close()


	// Auth service
	//auths := auth.NewStoreAuth(db, cfg.JWTKey)
	as, err := auth.NewAurum("http://localhost:8042")
	if err != nil {
		log.Fatal(err)
	}

	au := auth.NewAuthenticator(as, db)

	// Create initial user
	//initialUser(db, auths)

	// Executor
	var exec executor.Executor
	switch cfg.Executor {
	default:
		fallthrough
	case "docker":
		exec, err = executor.NewDockerExecutor(cfg.RunnerImage)
	case "kubernetes":
		exec, err = executor.NewKubernetesExecutor(cfg.RunnerImage, cfg.KubeNamespace, cfg.KubeConfigPath)
	}
	if err != nil {
		log.Fatalf("Starting %s executor failed: %v", cfg.Executor, err)
	}

	// Router
	rs := routes.New(cfg, db, au, exec)

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll)
	//r.Use(middleware.Compress(5))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", rs.HelloWorld)

	r.Post("/login", rs.Login)
	r.Get("/package", rs.GetPackages)

	// Protected Routes
	r.Group(func(r chi.Router) {
		// Veirfy jwt tokens
		r.Use(au.VerificationMiddleware)

		//r.Use(corsHandler)

		r.Post("/user", rs.AddUser)
		r.Delete("/user/{username}", rs.DeleteUser)
		r.Put("/user", rs.UpdateUser)
		r.Get("/users", rs.GetUsers)

		r.Post("/package", rs.AddPackage)

		r.Delete("/package/{pkg}", rs.RemovePackage)
		r.Put("/package/{pkg}", rs.UpdatePackage)

		r.Post("/package/{pkg}/upload", rs.UploadPackage)
		r.Put("/package/{pkg}/build", rs.TriggerBuild)
	})

	log.Fatal(http.ListenAndServe(cfg.Address, r))
}
