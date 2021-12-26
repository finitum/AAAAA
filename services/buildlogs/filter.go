package main

import (
	"github.com/finitum/AAAAA/pkg/models"
	"sort"
	"strconv"
	"strings"
)

// FilterJobs filters, sorts and paginates a list of jobs.
//
// * nameFilter is a keyword which if non-empty, must be included in the package name, or otherwise the job is filtered out.
// * statusFilter is a keyword which if non-empty, filters out any job that doesn't have a build status equal to the number
//   passed in. This parameter may start with a `!`, which negates the filter and filters out any job with a status equal to
//   the number passed in.
// * sortKey is the key in the remaining list of (filtered) jobs is sorted by before it's paginated. This parameter may
//   be either `time` or `name`. If this parameter isn't `time` or `name` it will automatically be sorted by time.
// * start and limit are used for pagination. Jobs returned are the filtered, sorted jobs sliced by [start:start+limit].
//   limit may be -1 to signify no limit.
func FilterJobs(jobs []models.Job, nameFilter, statusFilter, sortKey string, start, limit int) ([]models.Job, error) {
	if statusFilter != "" {
		reverse := false
		if statusFilter[0] == '!' {
			statusFilter = statusFilter[1:]
			reverse = true
		}

		statusNumber, err := strconv.Atoi(statusFilter)
		if err != nil {
			return nil, err
		}

		fc := 0
		for _, job := range jobs {
			if (job.Status == models.BuildStatus(statusNumber) && !reverse) ||
				(job.Status != models.BuildStatus(statusNumber) && reverse) {
				jobs[fc] = job
				fc++
			}
		}
		jobs = jobs[:fc]
	}

	if nameFilter != "" {
		fc := 0
		for _, job := range jobs {
			if strings.Contains(job.PackageName, nameFilter) {
				jobs[fc] = job
				fc++
			}
		}
		jobs = jobs[:fc]
	}

	sort.Slice(jobs, func(i, j int) bool {
		switch sortKey {
		case "name":
			name1 := jobs[i].PackageName
			name2 := jobs[j].PackageName
			return name1 < name2
		case "time":
			fallthrough
		default:
			return jobs[i].Time.After(jobs[j].Time)
		}
	})

	if len(jobs) > start {
		if limit == -1 {
			jobs = jobs[start:]
		} else {
			end := start + limit
			if end > len(jobs) {
				end = len(jobs)
			}

			jobs = jobs[start:end]
		}
	} else {
		jobs = []models.Job{}
	}

	return jobs, nil
}
