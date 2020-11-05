package store

import (
	"context"
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const gcTime = 5 * time.Minute

type Badger struct {
	db *badger.DB
}

func OpenBadger(path string) (*Badger, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, errors.Wrap(err, "opening badger store")
	}
	return &Badger{
		db,
	}, nil
}

func (b *Badger) StartGC(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(gcTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
			again:
				log.Debug("Garbage collection started")
				err := b.db.RunValueLogGC(0.7)
				if err == nil {
					goto again
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
