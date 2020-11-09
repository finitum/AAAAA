package main

import (
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testJobList = []models.Job{
	{
		PackageName: "ccc",
		Status:      0,
		Logs:        nil,
		Uuid:        "1",
		Time:        time.Unix(10, 10),
	},
	{
		PackageName: "aaa",
		Status:      1,
		Logs:        nil,
		Uuid:        "2",
		Time:        time.Unix(11, 10),
	},
	{
		PackageName: "bbb",
		Status:      0,
		Logs:        nil,
		Uuid:        "2",
		Time:        time.Unix(12, 10),
	},
}

func TestFilterLimitOne(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// By defaul it sorts based on time. So it should return the latest time (bbb)
	jobs, err := FilterJobs(jobs, "", "", "", "", "1")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].PackageName, "bbb")
}

func TestFilterLimitTwo(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// By defaul it sorts based on time. So it should return the latest two times (bbb, aaa)
	jobs, err := FilterJobs(jobs, "", "", "", "", "2")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 2)
	assert.Equal(t, jobs[0].PackageName, "bbb")
	assert.Equal(t, jobs[1].PackageName, "aaa")
}

func TestFilterLimitZero(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// Limit 0 should return nothing
	jobs, err := FilterJobs(jobs, "", "", "", "", "0")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 0)
}

func TestFilterStatusOne(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// There's only one with status 1 (bbb)
	jobs, err := FilterJobs(jobs, "", "1", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].PackageName, "aaa")
}

func TestFilterStatusZero(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// There are two  with status 0 (bbb and ccc) but ccc has the lowest time so should come last
	jobs, err := FilterJobs(jobs, "", "0", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 2)
	assert.Equal(t, jobs[0].PackageName, "bbb")
	assert.Equal(t, jobs[1].PackageName, "ccc")
}

func TestFilterStatusNotZero(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// There are two with status 0 (bbb and ccc) so !0 should return a
	jobs, err := FilterJobs(jobs, "", "!0", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].PackageName, "aaa")
}

func TestFilterNameExact(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// There are two with status 0 (bbb and ccc) so !0 should return a
	jobs, err := FilterJobs(jobs, "aaa", "", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].PackageName, "aaa")
}

func TestFilterNamePartial(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// There are two with status 0 (bbb and ccc) so !0 should return a
	jobs, err := FilterJobs(jobs, "a", "", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 1)
	assert.Equal(t, jobs[0].PackageName, "aaa")
}


func TestFilterSortTime(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// Sort based on time explicitly
	jobs, err := FilterJobs(jobs, "", "", "", "time", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 3)
	assert.Equal(t, jobs[0].PackageName, "bbb")
	assert.Equal(t, jobs[1].PackageName, "aaa")
	assert.Equal(t, jobs[2].PackageName, "ccc")
}

func TestFilterSortNothing(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// Sorting on nothing also sorts on time by default
	jobs, err := FilterJobs(jobs, "", "", "", "", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 3)
	assert.Equal(t, jobs[0].PackageName, "bbb")
	assert.Equal(t, jobs[1].PackageName, "aaa")
	assert.Equal(t, jobs[2].PackageName, "ccc")
}

func TestFilterSortName(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	jobs, err := FilterJobs(jobs, "", "", "", "name", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 3)
	assert.Equal(t, jobs[0].PackageName, "aaa")
	assert.Equal(t, jobs[1].PackageName, "bbb")
	assert.Equal(t, jobs[2].PackageName, "ccc")
}

func TestFilterStartZero(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// Starting at zero is the default
	jobs, err := FilterJobs(jobs, "", "", "0", "name", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 3)
	assert.Equal(t, jobs[0].PackageName, "aaa")
	assert.Equal(t, jobs[1].PackageName, "bbb")
	assert.Equal(t, jobs[2].PackageName, "ccc")
}

func TestFilterStartOne(t *testing.T) {
	jobs := make([]models.Job, len(testJobList))
	copy(jobs, testJobList)

	// Starting at zero is the default
	jobs, err := FilterJobs(jobs, "", "", "1", "name", "")
	assert.NoError(t, err)
	assert.Equal(t, len(jobs), 2)
	assert.Equal(t, jobs[0].PackageName, "bbb")
	assert.Equal(t, jobs[1].PackageName, "ccc")
}
