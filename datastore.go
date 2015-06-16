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
	if box.ds == nil {
		panic(ErrNoContext)
	}
	return box.ds
}

func AddDatastore(ctx context.Context, ds Datastore) context.Context {
	box := c(ctx)
	box.ds = ds
	return setC(ctx, box)
}
