package hfae

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
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
	keyString := key.Encode()
	_, err := mc.Codec.Get(ctx, keyString, dst)
	return err
}

func (mc *cache) Set(ctx context.Context, key *datastore.Key, dst interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			log.Debugf(ctx, "%s", err)
			log.Debugf(ctx, "%+v %+v", key, dst)
		}
	}()
	if key == nil {
		return datastore.ErrInvalidKey
	}

	if dst == nil {
		return datastore.ErrInvalidEntityType
	}

	err := mc.Codec.Set(ctx, &memcache.Item{
		Key:    key.Encode(),
		Object: dst,
	})

	return err
}

func (mc *cache) SetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	value := reflect.ValueOf(dst)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	dstLength := value.Len()
	if dstLength != len(keys) {
		log.Errorf(ctx, "SetMulti -> Expected: %d Found: %d", len(keys), dstLength)
		return nil
	}
	for i := 0; i < dstLength; i++ {
		if ctx.Err() != nil {
			i = dstLength
		} else {
			item := value.Index(i)
			if item.Interface() != nil {
				mc.Set(ctx, keys[i], item.Interface())
			}
		}
	}

	return nil
}
