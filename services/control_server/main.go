package main

import (
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := store.OpenBadgerStore(os.TempDir() + "/AAAAA")
	if err != nil {
		log.Fatal(err)
	}

	rs := NewRoutes(db)

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Compress(5))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", rs.HelloWorld)

	log.Fatal(http.ListenAndServe(":5000", r))
}
