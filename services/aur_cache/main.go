package main

import (
	"context"
	"github.com/finitum/AAAAA/internal/cors"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
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

	cachetype := os.Getenv("CACHE_TYPE")
	var cache store.Cache

	switch strings.ToLower(cachetype) {
	default:
		fallthrough
	case "ristretto":
		ristretto, err := store.NewRistretto()
		if err != nil {
			log.Fatalf("Couldn't open ristretto cache: %v", err)
		}
		defer ristretto.Close()
		cache = ristretto
	case "badger":
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		badger, err := store.OpenBadger(os.TempDir() + "/AAAAA-cache")
		if err != nil {
			log.Fatalf("Couldn't open ristretto cache: %v", err)
		}
		defer badger.Close()
		badger.StartGC(ctx)
		cache = badger
	}

	r.Get("/search/{term}", search(cache))
	r.Get("/info/{name}", info(cache))

	log.Fatal(http.ListenAndServe(":5001", r))
}

const maxReturn int = 50

func info(cache store.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		entry, err := cache.GetInfoEntry(name)
		if err == nil {
			_ = render.Render(w, r, entry)
			return
		}

		if err != store.ErrNotExists {
			log.Errorf("Unexpected error when retrieving cached info entry: %v", err)
		} else {
			log.Trace("Info cache miss")
		}

		res, err := aur.SendInfoRequest(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if len(res.Results) < 1 {
			http.Error(w, "no results", http.StatusBadGateway)
			return
		}

		if err := cache.SetInfoEntry(name, &res.Results[0]); err != nil {
			log.Error("saving to cache failed")
		}

		_ = render.Render(w, r, &res.Results[0])
	}
}

func search(cache store.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		term := chi.URLParam(r, "term")
		if len(term) < 3 {
			http.Error(w, "term too short", http.StatusBadRequest)
			return
		}

		cachedResult, exact, err := store.GetPartialCacheEntry(cache, term)
		if err == nil {
			log.Trace("Cache hit!")
			if exact {
				if len(cachedResult) > maxReturn {
					cachedResult = cachedResult[:maxReturn]
				}
			} else {
				fi := 0
				// The cached result may contain too many entries. Manual filter required
				for _, item := range cachedResult {
					if strings.Contains(item.Description, term) || strings.Contains(item.Name, term) {
						cachedResult[fi] = item
						fi++
					}

					if fi > maxReturn {
						break
					}
				}

				cachedResult = cachedResult[:fi]
			}

			_ = render.Render(w, r, cachedResult)
			return
		}

		// otherwise it's just a cache miss
		if err != store.ErrNotExists {
			log.Errorf("An unexpected error occurred while retrieving cache results. Attempting a non-cached lookup (%v)", err)
		} else {
			log.Trace("Cache miss!")
		}

		res, err := aur.SendResultsRequest(term)
		if err != nil {
			http.Error(w, "received error from aur rpc", http.StatusBadGateway)
			return
		}

		// Sort first
		res.Results.SortByPopularity()

		// Store later
		err = cache.SetResultsEntry(term, res.Results)
		if err != nil {
			// A cache error means we can still return the results we looked up.
			log.Errorf("Failed to set cache entry (%v)", err)
		}

		if len(res.Results) > maxReturn {
			res.Results = res.Results[:maxReturn]
		}
		_ = render.Render(w, r, res.Results)
	}
}
