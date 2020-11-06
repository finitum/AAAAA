package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))
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

	r.Get("/search/{term}", proxy(cache))

	log.Fatal(http.ListenAndServe(":5001", r))
}

const aurRpcQueryFmt = "https://aur.archlinux.org/rpc/?v=5&type=search&arg=%s"
const maxReturn int = 50

func proxy(cache store.Cache) http.HandlerFunc {
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

		resp, err := http.Get(fmt.Sprintf(aurRpcQueryFmt, term))
		if err != nil {
			http.Error(w, "received error from aur rpc", http.StatusBadGateway)
			return
		}

		var res aur.ExtendedResults
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			http.Error(w, "couldn't decode result", http.StatusBadGateway)
			log.Error(err)
			return
		}

		if res.Error != "" {
			http.Error(w, res.Error, http.StatusBadRequest)
			return
		}

		// Sort first
		res.Results.SortByPopularity()

		// Store later
		err = cache.SetEntry(term, res.Results)
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
