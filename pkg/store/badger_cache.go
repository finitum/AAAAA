package store

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v2"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/pkg/errors"
	"time"
)

type BadgerCache struct {
	db *badger.DB
}

const cachePrefix = "cache_"
const cacheTTL = 30 * time.Minute

func OpenBadgerCache(path string) (*BadgerCache, error) {
	db, err := badger.Open(badger.DefaultOptions(path + ".cache"))
	if err != nil {
		return nil, errors.Wrap(err, "opening badger store")
	}

	// TODO: Schedule garbage collection

	return &BadgerCache{
		db,
	}, nil
}

func (b *BadgerCache) SetEntry(searchterm string, result aur.Results) error {
	return b.db.Update(func(txn *badger.Txn) error {
		var value bytes.Buffer

		enc := gob.NewEncoder(&value)
		err := enc.Encode(result)
		if err != nil {
			return errors.Wrap(err, "gob encode")
		}

		// Add the main key to the store
		mainEntryKey := []byte(cachePrefix + searchterm)
		mainEntry := badger.NewEntry(mainEntryKey, value.Bytes()).WithTTL(cacheTTL)
		err = txn.SetEntry(mainEntry)
		if err != nil {
			return errors.Wrap(err, "badger main transaction")
		}

		return nil
	})
}

func (b *BadgerCache) GetEntry(searchterm string) (aur.Results, bool, error) {
	var result aur.Results
	exact := true

	err := b.db.View(func(txn *badger.Txn) error {
		// Add all the smaller strings to the database as well, to ensure a quicker lookup
		for i := len(searchterm); i > 2; i-- {
			actualSearchterm := []byte(cachePrefix + searchterm[:i])

			// Get a value from the store. This may be a pointer to another key
			item, err := txn.Get(actualSearchterm)
			if err == badger.ErrKeyNotFound {
				// After the first iteration it's not exact anymore
				exact = false
				continue
			} else if err != nil {
				return errors.Wrap(err, "badger get")
			}

			err = item.Value(func(val []byte) error {
				buf := bytes.NewBuffer(val)

				dec := gob.NewDecoder(buf)
				return errors.Wrap(dec.Decode(&result), "gob decode")
			})

			return errors.Wrap(err, "badger read")
		}

		return ErrNotExists
	})

	if err != nil {
		return nil, false, err
	}

	return result, exact, err
}
