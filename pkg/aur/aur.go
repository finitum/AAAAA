package aur

import (
	"net/http"
	"sort"
)

type Result struct {
	ID             int
	Name           string
	PackageBaseID  int
	PackageBase    string
	Version        string
	Description    string
	URL            string
	NumVotes       int
	Popularity     float64
	OutOfDate      int
	Maintainer     string
	FirstSubmitted int
	LastModified   int
	URLPath        string
}

type Results []Result

func (r Results) SortByPopularity() {
	sort.Slice(r, func(i, j int) bool {
		return r[i].Popularity > r[j].Popularity
	})
}

func (r Results) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type ExtendedResults struct {
	Version     int     `json:"version"`
	Type        string  `json:"type"`
	Results     Results `json:"results"`
	ResultCount int     `json:"resultcount"`
	Error       string  `json:"error"`
}
