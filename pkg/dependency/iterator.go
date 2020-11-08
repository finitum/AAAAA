package dependency

import (
	"github.com/Workiva/go-datastructures/queue"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Iterator struct {
	q    *queue.Queue
	next string
}

func NewIterator(deps []string) (*Iterator, error) {
	i := &Iterator{q: queue.New(int64(len(deps)))}

	for _, dep := range deps {
		if err := i.q.Put(dep); err != nil {
			return nil, errors.Wrap(err, "unable to add dependencies to queue")
		}
	}

	return i, nil
}

func (it *Iterator) Next() bool {
	if it.q.Empty() {
		return false
	}

	res, err := it.q.Get(1)
	if err != nil {
		// Shouldn't happen as Next() should always be called before Item() and never after Close(),
		// but logging the error just in case.
		log.Errorf("error while retrieving next item from iterator queue: %v", err)
		return false
	}

	// Should be safe, as this struct is the only thing directly accessing the queue
	it.next = res[0].(string)
	return true
}

func (it *Iterator) Item() string {
	return it.next
}

func (it *Iterator) Push(dep string) {
	err := it.q.Put(dep)
	if err != nil {
		// Shouldn't happen as Push() should never be called after Close().
		log.Errorf("error while pushing item to iterator queue: %v", err)
	}
}

func (it *Iterator) Close() {
	it.q.Dispose()
}
