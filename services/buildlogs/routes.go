package main

import (
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/finitum/AAAAA/pkg/store"
	"github.com/finitum/AAAAA/services/control_server/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Routes struct {
	jobs store.JobStore
}

func NewRoutes(j store.JobStore) *Routes {
	return &Routes{
		j,
	}
}

func (rs *Routes) NewJob(w http.ResponseWriter, r *http.Request) {
	pkgname := chi.URLParam(r, "pkgname")

	job, err := rs.jobs.NewJob(pkgname)
	if err != nil {
		_ = render.Render(w, r, routes.ErrServerError(err))
		log.Errorf("failed to create new job (%v)", err)
		return
	}

	_ = render.Render(w, r, job)
	w.WriteHeader(http.StatusCreated)
}

func (rs *Routes) GetJob(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	job, err := rs.jobs.GetJob(uuid)
	if err != nil {
		_ = render.Render(w, r, routes.ErrServerError(err))
		log.Errorf("failed to get job (%v)", err)
		return
	}

	_ = render.Render(w, r, job)
}

func (rs *Routes) GetLogs(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	dbLogs, err := rs.jobs.GetLogs(uuid)
	if err != nil {
		_ = render.Render(w, r, routes.ErrServerError(err))
		log.Errorf("failed to get logs (%v)", err)
		return
	}

	logs := make([]render.Renderer, len(dbLogs))
	for i, logLine := range logs {
		logs[i] = logLine
	}

	_ = render.RenderList(w, r, logs)
}

// GetJobs possible routes:
//	All jobs: /jobs
// 	Get 10 jobs: /jobs?limit=10
// 	Get jobs starting at job 10: /jobs?start=10
// 	Get only the 10th job: /jobs?start=10&limit=1
//
// # Sorting is performed before filtering and limiting
// 	Sort jobs by time: /jobs?sort=time
// 	Sort jobs by package name: /jobs?sort=name
//
// # Filtering is performed before sorting
//  Return only jobs with status 0 (pending) : /jobs?status=0
//  Return only jobs with a status that's not 0 (pending) : /jobs?status=!0
//  Return only jobs with `aaa` in the name: /jobs?name=aaa
func (rs *Routes) GetJobs(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	start := r.URL.Query().Get("start")
	sortKey := r.URL.Query().Get("sort")
	statusFilter := r.URL.Query().Get("status")
	nameFilter := r.URL.Query().Get("name")

	var err error

	startNum := 0
	if start != "" {
		startNum, err = strconv.Atoi(start)
		if err != nil {
			_ = render.Render(w, r, routes.ErrInvalidRequest(err))
			log.Errorf("Couldn't convert status to number (%v)", err)
			return
		}
		if startNum < 0 {
			startNum = 0
		}
	}

	limitNum := -1
	if limit != "" {
		limitNum, err = strconv.Atoi(limit)
		if err != nil {
			_ = render.Render(w, r, routes.ErrInvalidRequest(err))
			log.Errorf("Couldn't convert status to number (%v)", err)
			return
		}
		if limitNum < 0 {
			limitNum = -1
		}
	}

	// We will never return more than 5000 results
	if limitNum > 5000 {
		limitNum = 5000
	}

	dbJobs, err := rs.jobs.GetJobs()
	if err != nil {
		// TODO: maybe put these error functions in some kind of shared module (internal maybe?)
		_ = render.Render(w, r, routes.ErrServerError(err))
		log.Errorf("failed to get jobs (%v)", err)
		return
	}

	dbJobs, err = FilterJobs(dbJobs, nameFilter, statusFilter, sortKey, startNum, limitNum)
	if err != nil {
		_ = render.Render(w, r, routes.ErrInvalidRequest(err))
		log.Errorf("Couldn't convert status to number (%v)", err)
		return
	}

	jobs := make([]render.Renderer, len(dbJobs))
	for i, logLine := range dbJobs {
		jobs[i] = logLine
	}

	_ = render.RenderList(w, r, jobs)
}

func (rs *Routes) AddLogs(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	var logLine models.LogLine

	if err := render.Bind(r, &logLine); err != nil {
		_ = render.Render(w, r, routes.ErrInvalidRequest(err))
		return
	}

	err := rs.jobs.AppendToJobLog(uuid, &logLine)
	if err != nil {
		_ = render.Render(w, r, routes.ErrServerError(err))
		log.Errorf("failed to add to logs (%v)", err)
		return
	}
}
