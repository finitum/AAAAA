package main

import (
	"encoding/json"
	"fmt"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main()  {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/search/{term}", proxy)

	log.Fatal(http.ListenAndServe(":5001", r))
}

const aurRpcQueryFmt = "https://aur.archlinux.org/rpc/?v=5&type=search&arg=%s"
const maxReturn int = 50

func proxy(w http.ResponseWriter, r *http.Request) {
	term := chi.URLParam(r, "term")
	if len(term) < 3 {
		http.Error(w, "term too short", http.StatusBadRequest)
		return
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

	// cache
	res.Results.SortByPopularity()


	if len(res.Results) > maxReturn {
		res.Results = res.Results[:maxReturn]
	}

	_ = render.Render(w, r, res.Results)
}
