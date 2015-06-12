package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func ds(ctx context.Context) *db {
	return &db{}
}

type db struct{}

func (ds *db) Delete(ctx context.Context, key *datastore.Key) error {
	return datastore.Delete(ctx, key)
}

func (ds *db) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	return datastore.DeleteMulti(ctx, keys)
}

func (ds *db) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	return datastore.Get(ctx, key, dst)
}

func (ds *db) GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	return datastore.GetMulti(ctx, keys, dst)
}

func (ds *db) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	keyResult, err := datastore.Put(ctx, key, src)
	return keyResult, err
}

func (ds *db) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	keyResult, err := datastore.PutMulti(ctx, keys, src)
	return keyResult, err
}
