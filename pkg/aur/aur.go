package aur

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

const aurInfoQuery = "https://aur.archlinux.org/rpc/?v=5&type=info&arg=%s"
const aurSearchQuery = "https://aur.archlinux.org/rpc/?v=5&type=search&arg=%s"

type SearchResult struct {
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

type InfoResult struct {
	SearchResult
	Depends     []string
	MakeDepends []string
	License     []string
	Keywords    []string
}

func (i *InfoResult) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

type Results []SearchResult

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

type ExtendedInfoResults struct {
	Version     int          `json:"version"`
	Type        string       `json:"type"`
	ResultCount int          `json:"resultcount"`
	Results     []InfoResult `json:"results"`
	Error       string       `json:"error"`
}

func SendInfoRequest(pkg string) (res ExtendedInfoResults, err error) {
	resp, err := http.Get(fmt.Sprintf(aurInfoQuery, pkg))
	if err != nil {
		return res, errors.Wrap(err, "received error from aur rpc")
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Error(err)
		return res, errors.Wrap(err, "couldn't decode result")
	}

	if res.Error != "" {
		return res, errors.New(fmt.Sprintf("error from aur: %v", res.Error))
	}

	return
}

func SendResultsRequest(term string) (res ExtendedResults, _ error) {
	resp, err := http.Get(fmt.Sprintf(aurSearchQuery, term))
	if err != nil {
		return res, errors.Wrap(err, "received error from aur rpc")
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return res, errors.Wrap(err, "couldn't decode result")
	}

	if res.Error != "" {
		return res, errors.New(fmt.Sprintf("error from aur: %v", res.Error))
	}

	return
}
