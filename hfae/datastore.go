package hfae

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type store struct{}

func (ds *store) Delete(ctx context.Context, key *datastore.Key) error {
	return datastore.Delete(ctx, key)
}

func (ds *store) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return datastore.DeleteMulti(ctx, keys)
}

func (ds *store) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	return datastore.Get(ctx, key, dst)
}

func (ds *store) GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	return datastore.GetMulti(ctx, keys, dst)
}

func (ds *store) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	keyResult, err := datastore.Put(ctx, key, src)
	return keyResult, err
}

func (ds *store) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	keyResult, err := datastore.PutMulti(ctx, keys, src)
	return keyResult, err
}
