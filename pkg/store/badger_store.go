package store

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v2"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/pkg/errors"
)

const pkgPrefix = "pkg_"
const userPrefix = "user_"

func (b *Badger) GetPackage(name string) (*models.Pkg, error) {
	var pkg *models.Pkg

	return pkg, b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(pkgPrefix + name))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}

		return errors.Wrap(item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&pkg), "gob decode")
		}), "badger read")
	})
}

func (b *Badger) AddPackage(pkg *models.Pkg) error {
	return b.db.Update(func(txn *badger.Txn) error {
		var value bytes.Buffer

		enc := gob.NewEncoder(&value)
		err := enc.Encode(pkg)
		if err != nil {
			return errors.Wrap(err, "gob encode")
		}

		return errors.Wrap(txn.Set([]byte(pkgPrefix+pkg.Name), value.Bytes()), "badger transaction")
	})
}

func (b *Badger) DelPackage(name string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return errors.Wrap(txn.Delete([]byte(pkgPrefix+name)), "badger transaction")
	})
}

func (b *Badger) AllPackages() (pkgs []*models.Pkg, err error) {
	return pkgs, b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(pkgPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var pkg models.Pkg
			err := item.Value(func(val []byte) error {
				buf := bytes.NewBuffer(val)

				dec := gob.NewDecoder(buf)
				return errors.Wrap(dec.Decode(&pkg), "gob decode")
			})
			pkgs = append(pkgs, &pkg)
			if err != nil {
				return errors.Wrap(err, "badger iteration")
			}
		}
		return nil
	})
}

func (b *Badger) AllPackageNames() (names []string, _ error) {
	return b.allKeysWithPrefix([]byte(pkgPrefix))
}

func (b *Badger) GetUser(name string) (user *models.User, err error) {
	err = b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userPrefix + name))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}

		return errors.Wrap(item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&user), "gob decode")
		}), "badger read")
	})
	return
}

func (b *Badger) AddUser(user *models.User) error {
	return b.db.Update(func(txn *badger.Txn) error {
		var value bytes.Buffer

		enc := gob.NewEncoder(&value)
		err := enc.Encode(user)
		if err != nil {
			return errors.Wrap(err, "gob encode")
		}

		return errors.Wrap(txn.Set([]byte(userPrefix+user.Username), value.Bytes()), "badger transaction")
	})
}

func (b *Badger) DelUser(username string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return errors.Wrap(txn.Delete([]byte(userPrefix+username)), "badger transaction")
	})
}

func (b *Badger) AllUserNames() (users []string, _ error) {
	return b.allKeysWithPrefix([]byte(userPrefix))
}

func (b *Badger) AllUsers() (users []*models.User, _ error) {
	return users, b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte(userPrefix)); it.ValidForPrefix([]byte(userPrefix)); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				buf := bytes.NewBuffer(val)
				dec := gob.NewDecoder(buf)

				var user models.User
				err := dec.Decode(&user)
				if err != nil {
					return errors.Wrap(err, "gob decode")
				}
				users = append(users, &user)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return errors.Wrap(nil, "wait what")
	})
}

func (b *Badger) allKeysWithPrefix(prefix []byte) (names []string, _ error) {
	return names, b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()

			names = append(names, string(k[len(prefix):]))
		}

		return errors.Wrap(nil, "wait what")
	})
}
