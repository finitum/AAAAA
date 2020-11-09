package store

import (
	"bytes"
	"encoding/gob"
	"github.com/dgraph-io/badger/v2"
	"github.com/finitum/AAAAA/pkg/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type JobStoreWrapper struct {
	*Badger

	sync.Mutex
	callbacks map[string][]func(line *models.LogLine)
}

func NewJobStore(badger *Badger) JobStore {
	return &JobStoreWrapper{
		Badger:    badger,
		callbacks: make(map[string][]func(line *models.LogLine)),
	}
}

func (b *JobStoreWrapper) NewJob(name string) (*models.Job, error) {
	jid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "uuid")
	}

	job := models.Job{
		PackageName: name,
		Status:      models.BuildStatusPending,
		Logs:        nil,
		Uuid:        jid.String(),
		Time:        time.Now(),
	}

	err = b.db.Update(func(txn *badger.Txn) error {
		var value bytes.Buffer

		enc := gob.NewEncoder(&value)
		err := enc.Encode(job)
		if err != nil {
			return errors.Wrap(err, "gob encode")
		}

		entry := badger.NewEntry([]byte(jobPrefix+jid.String()), value.Bytes()).WithTTL(jobTTL)
		return errors.Wrap(txn.SetEntry(entry), "badger transaction")
	})

	return &job, err
}

func (b *JobStoreWrapper) AppendToJobLog(jid string, l *models.LogLine) error {
	for _, cb := range b.callbacks[jid] {
		cb(l)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		var job models.Job

		// Get the job
		item, err := txn.Get([]byte(jobPrefix + jid))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}
		err = item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&job), "gob decode")
		})
		if err != nil {
			return err
		}

		// Update the job
		job.Logs = append(job.Logs, *l)

		// Put the job back
		var value bytes.Buffer
		enc := gob.NewEncoder(&value)
		err = enc.Encode(job)
		if err != nil {
			return errors.Wrap(err, "gob encode")
		}

		entry := badger.NewEntry([]byte(jobPrefix+jid), value.Bytes()).WithTTL(jobTTL)
		return errors.Wrap(txn.SetEntry(entry), "badger transaction")
	})
}

func (b *JobStoreWrapper) SetJobStatus(jid string, status models.BuildStatus) error {
	return b.db.Update(func(txn *badger.Txn) error {
		var job models.Job

		// Get the job
		item, err := txn.Get([]byte(jobPrefix + jid))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}
		err = item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&job), "gob decode")
		})
		if err != nil {
			return err
		}

		// Update the job
		job.Status = status

		// Put the job back
		var value bytes.Buffer
		enc := gob.NewEncoder(&value)
		err = enc.Encode(job)
		if err != nil {
			return errors.Wrap(err, "gob encode")
		}

		entry := badger.NewEntry([]byte(jobPrefix+jid), value.Bytes()).WithTTL(jobTTL)
		return errors.Wrap(txn.SetEntry(entry), "badger transaction")
	})
}

func (b *JobStoreWrapper) GetLogs(jid string) (logs []models.LogLine, _ error) {
	return logs, b.db.View(func(txn *badger.Txn) error {
		var job models.Job

		// Get the job
		item, err := txn.Get([]byte(jobPrefix + jid))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}
		err = item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&job), "gob decode")
		})
		if err != nil {
			return err
		}

		logs = job.Logs

		return nil
	})
}

func (b *JobStoreWrapper) GetJobs() (jobs []models.Job, _ error) {
	return jobs, b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(jobPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var job models.Job
			err := item.Value(func(val []byte) error {
				buf := bytes.NewBuffer(val)

				dec := gob.NewDecoder(buf)
				return errors.Wrap(dec.Decode(&job), "gob decode")
			})
			jobs = append(jobs, job)
			if err != nil {
				return errors.Wrap(err, "badger iteration")
			}
		}
		return nil
	})
}

func (b *JobStoreWrapper) AddLogListener(uuid string, cb func(line *models.LogLine)) {
	b.Lock()
	defer b.Unlock()

	b.callbacks[uuid] = append(b.callbacks[uuid], cb)
}

func (b *JobStoreWrapper) GetJob(jid string) (*models.Job, error) {
	var job models.Job

	return &job, b.db.View(func(txn *badger.Txn) error {

		// Get the job
		item, err := txn.Get([]byte(jobPrefix + jid))
		if err == badger.ErrKeyNotFound {
			return ErrNotExists
		} else if err != nil {
			return errors.Wrap(err, "badger get")
		}
		err = item.Value(func(val []byte) error {
			buf := bytes.NewBuffer(val)

			dec := gob.NewDecoder(buf)
			return errors.Wrap(dec.Decode(&job), "gob decode")
		})
		if err != nil {
			return err
		}

		return nil
	})
}
