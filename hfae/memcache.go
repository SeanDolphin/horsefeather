package hfae

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

type cache struct {
	Codec memcache.Codec
}

func (mc *cache) Delete(ctx context.Context, key *datastore.Key) error {
	return memcache.Delete(ctx, key.Encode())
}

func (mc *cache) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	simpleKeys := make([]string, len(keys))
	for i, key := range keys {
		simpleKeys[i] = key.Encode()
	}

	return memcache.DeleteMulti(ctx, simpleKeys)
}

func (mc *cache) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	_, err := mc.Codec.Get(ctx, key.Encode(), dst)
	return err
}

func (mc *cache) Set(ctx context.Context, key *datastore.Key, dst interface{}) error {
	return mc.Codec.Set(ctx, &memcache.Item{
		Key:    key.Encode(),
		Object: dst,
	})
}

func (mc *cache) SetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	value := reflect.ValueOf(dst)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		mc.Set(ctx, keys[i], item.Interface())
	}

	return nil
}
