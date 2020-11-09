package main

import (
	"github.com/finitum/AAAAA/pkg/models"
	"sort"
	"strconv"
	"strings"
)

func FilterJobs(jobs []models.Job, nameFilter, statusFilter, start, sortKey, limit string) ([]models.Job, error) {
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

	if start != "" {
		startNum, err := strconv.Atoi(start)
		if err != nil {
			return nil, err
		}

		if len(jobs) > startNum {
			jobs = jobs[startNum:]
		} else {
			jobs = []models.Job{}
		}
	}

	limitNum := 5000
	var err error

	if limit != "" {
		limitNum, err = strconv.Atoi(limit)

		if err != nil {
			return nil, err
		}
	}

	if len(jobs) > limitNum {
		jobs = jobs[:limitNum]
	}

	return jobs, nil
}
