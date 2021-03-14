package store

import (
	"github.com/dgraph-io/ristretto"
	"github.com/pkg/errors"

	"github.com/finitum/AAAAA/pkg/aur"
)

type Ristretto struct {
	cache *ristretto.Cache
}

// NewRistretto creates a new Cache based on ristretto.
func NewRistretto() (*Ristretto, error) {
	// TODO: these are sane defaults, we may want to make them configurable
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating ristretto cache")
	}

	return &Ristretto{cache}, nil
}

func (r *Ristretto) SetResultsEntry(term string, result aur.Results) error {
	r.cache.SetWithTTL(resultsPrefix+term, result, int64(len(result)), cacheTTL)
	return nil
}

func (r *Ristretto) GetResultsEntry(term string) (aur.Results, error) {
	value, found := r.cache.Get(resultsPrefix + term)
	if !found {
		return nil, ErrNotExists
	}

	return value.(aur.Results), nil
}

func (r *Ristretto) SetInfoEntry(name string, result *aur.InfoResult) error {
	r.cache.SetWithTTL(infoPrefix+name, result, 1, cacheTTL)
	return nil
}

func (r *Ristretto) GetInfoEntry(name string) (*aur.InfoResult, error) {
	value, found := r.cache.Get(infoPrefix + name)
	if !found {
		return nil, ErrNotExists
	}

	return value.(*aur.InfoResult), nil
}

func (r *Ristretto) Close() {
	r.cache.Close()
}
