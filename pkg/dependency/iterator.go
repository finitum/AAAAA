package dependency

import (
	"github.com/Workiva/go-datastructures/queue"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Iterators represents an iterator which has the option to iterate over all dependencies, and also allows re-adding
// a Dependency to the back of the iterator queue.
type Iterator struct {
	q    *queue.Queue
	next Dependency
}

// NewIterator returns a new Iterator based on the given Dependency slice.
func NewIterator(deps []Dependency) (*Iterator, error) {
	i := &Iterator{q: queue.New(int64(len(deps)))}

	for _, dep := range deps {
		if err := i.q.Put(dep); err != nil {
			return nil, errors.Wrap(err, "unable to add Dependencies to queue")
		}
	}

	return i, nil
}

// Next returns true if there is a value to consume, false if not. Calling next multiple times without calling Item()
// will silently consume items.
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
	it.next, _ = res[0].(Dependency)
	return true
}

// Item returns the current Dependency. Calling Item() multiple times without calling Next(), will return the same
// object multiple times.
func (it *Iterator) Item() Dependency {
	return it.next
}

// Push re-adds a Dependency to the iterator queue.
func (it *Iterator) Push(dep Dependency) {
	err := it.q.Put(dep)
	if err != nil {
		// Shouldn't happen as Push() should never be called after Close().
		log.Errorf("error while pushing item to iterator queue: %v", err)
	}
}

// Close disposes of internal datastructures. Calling Next(), Item(), or Push() after Close() is undefined behaviour.
func (it *Iterator) Close() {
	it.q.Dispose()
}
