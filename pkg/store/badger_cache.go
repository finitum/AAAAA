package store

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v2"
	"github.com/finitum/AAAAA/pkg/aur"
	"github.com/pkg/errors"
)

const cachePrefix = "cache_"

func (b *Badger) SetEntry(searchterm string, result aur.Results) error {
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

func (b *Badger) GetEntry(term string) (result aur.Results, _ error) {
	return result, b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(cachePrefix + term))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}

		return item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&result), "gob decode")
		})
	})
}
