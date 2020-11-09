package main

import (
	"context"
	"github.com/finitum/AAAAA/internal/cors"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	badger, err := store.OpenBadger(os.TempDir() + "/AAAAA-jobs")
	if err != nil {
		log.Fatalf("Couldn't open ristretto cache: %v", err)
	}
	defer badger.Close()
	badger.StartGC(ctx)

	js := store.NewJobStore(badger)
	rs := NewRoutes(js)

	r.Post("/job/{pkgname}", rs.NewJob)
	r.Get("/job/{uuid}", rs.GetJob)
	r.Get("/jobs", rs.GetJobs)
	r.Get("/job/{uuid}/logs", rs.GetLogs)
	r.Post("/job/{uuid}/logs", rs.AddLogs)

	log.Fatal(http.ListenAndServe(":5002", r))
}
