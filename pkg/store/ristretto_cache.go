package store

import (
	"github.com/dgraph-io/ristretto"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/pkg/errors"
)

type Ristretto struct {
	cache *ristretto.Cache
}

// NewRistretto creates a new Cache based on ristretto
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

func (r *Ristretto) SetEntry(term string, result aur.Results) error {
	r.cache.SetWithTTL(term, result, int64(len(result)), cacheTTL)
	return nil
}

func (r *Ristretto) GetEntry(term string) (aur.Results, error) {
	value, found := r.cache.Get(term)
	if !found {
		return nil, ErrNotExists
	}

	return value.(aur.Results), nil
}

func (r *Ristretto) Close() {
	r.cache.Close()
}
