package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Memcache interface {
	Delete(context.Context, *datastore.Key) error
	DeleteMulti(context.Context, []*datastore.Key) error
	Get(context.Context, *datastore.Key, interface{}) error
	Set(context.Context, *datastore.Key, interface{}) error
	SetMulti(context.Context, []*datastore.Key, interface{}) error
}

func mc(ctx context.Context) Memcache {
	box := c(ctx)
	box.RLock()
	defer box.RUnlock()

	if box.mc == nil {
		panic(ErrNoContext)
	}
	return box.mc
}

func AddMemcache(ctx context.Context, mc Memcache) context.Context {
	box := c(ctx)
	box.Lock()
	defer box.Unlock()

	box.mc = mc
	return setC(ctx, box)
}

func OnlyMemcache(ctx context.Context, flag bool) context.Context {
	box := c(ctx)
	box.Lock()
	defer box.Unlock()

	box.noDS = flag
	return setC(ctx, box)
}

func IsMemcacheAllowed(ctx context.Context) bool {
	box := c(ctx)
	box.RLock()
	defer box.RUnlock()

	return !box.noMC
}
