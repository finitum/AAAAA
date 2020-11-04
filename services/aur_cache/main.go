package main

import (
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

func main()  {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	cache, err := store.OpenBadgerCache(os.TempDir() + "/AAAAA")
	if err != nil {
		log.Fatalf("Couldn't open store (%v)", err)
	}

	r.Get("/search/{term}", func(writer http.ResponseWriter, request *http.Request) {
		proxy(cache, writer, request)
	})

	log.Fatal(http.ListenAndServe(":5001", r))
}

const aurRpcQueryFmt = "https://aur.archlinux.org/rpc/?v=5&type=search&arg=%s"
const maxReturn int = 50

func proxy(cache store.Cache, w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	if len(term) < 3 {
		http.Error(w, "term too short", http.StatusBadRequest)
		return
	}

	cachedResult, err := cache.GetEntry(term)
	if err == nil {
		log.Info("Cache hit!")
		filteredResult := aur.Results(make([]aur.Result, len(cachedResult)))

		fi := 0
		// The cached result may contain too many entries. Manual filter required
		for _, i := range cachedResult {
			if strings.Contains(i.Description, term) || strings.Contains(i.Name, term) {
				filteredResult[fi] = i
				fi++
			}
		}

		_ = render.Render(w, r, filteredResult)
		return
	}

	// otherwise it's just a cache miss
	if err != store.ErrNotExists {
		log.Errorf("An unexpected error occurred while retrieving cache results. Attempting a non-cached lookup (%v)", err)
	} else {
		log.Info("Cache miss!")
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
