package horsefeather

import (
	"sync"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

type result struct {
	sync.RWMutex
	items map[*datastore.Key]interface{}
}

func newResult() *result {
	return &result{
		items: map[*datastore.Key]interface{}{},
	}
}

func (r *result) Set(ctx context.Context, key *datastore.Key, item interface{}) {
	r.Lock()
	defer r.Unlock()
	r.items[key] = item
}

func (r *result) Get(ctx context.Context, key *datastore.Key) interface{} {
	r.RLock()
	defer r.RUnlock()
	return r.items[key]
}

func (r *result) RemainingKeys() []*datastore.Key {
	r.RLock()
	defer r.RUnlock()
	remainingKeys := make([]*datastore.Key, 0, len(r.items))
	for key, item := range r.items {
		if item == nil {
			remainingKeys = append(remainingKeys, key)
		}
	}
	return remainingKeys
}
