package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Datastore interface {
	Delete(context.Context, *datastore.Key) error
	DeleteMulti(context.Context, []*datastore.Key) error
	Get(context.Context, *datastore.Key, interface{}) error
	GetMulti(context.Context, []*datastore.Key, interface{}) error
	Put(context.Context, *datastore.Key, interface{}) (*datastore.Key, error)
	PutMulti(context.Context, []*datastore.Key, interface{}) ([]*datastore.Key, error)
}

func ds(ctx context.Context) Datastore {
	box := c(ctx)
	box.RLock()
	defer box.RUnlock()

	if box.ds == nil {
		panic(ErrNoContext)
	}
	return box.ds
}

func AddDatastore(ctx context.Context, ds Datastore) context.Context {
	box := c(ctx)
	box.Lock()
	defer box.Unlock()

	box.ds = ds
	return setC(ctx, box)
}

func OnlyDatastore(ctx context.Context, flag bool) context.Context {
	box := c(ctx)
	box.Lock()
	defer box.Unlock()

	box.noMC = flag
	return setC(ctx, box)
}

func IsDatastoreAllowed(ctx context.Context) bool {
	box := c(ctx)
	box.RLock()
	defer box.RUnlock()

	return !box.noDS
}
