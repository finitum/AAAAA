package main

import (
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	r.Use(middleware.Logger)

	r.Get("/", rs.HelloWorld)

	log.Fatal(http.ListenAndServe(":5000", r))
}
