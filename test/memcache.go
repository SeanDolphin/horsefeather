package test

import (
	"encoding/json"
	"errors"
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func NewCache() *cache {
	return &cache{
		items: map[string][]byte{},
	}
}

type cache struct {
	items map[string][]byte
}

func (mc *cache) Delete(ctx context.Context, key *datastore.Key) error {
	delete(mc.items, key.Encode())
	return nil
}

func (mc *cache) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	for _, key := range keys {
		mc.Delete(ctx, key)
	}
	return nil
}

func (mc *cache) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	data, ok := mc.items[key.Encode()]
	if ok {
		return json.Unmarshal(data, &dst)
	}
	return errors.New("no item")
}

func (mc *cache) Set(ctx context.Context, key *datastore.Key, src interface{}) error {
	data, err := json.Marshal(src)
	if err == nil {
		mc.items[key.Encode()] = data
	}
	return err
}

func (mc *cache) SetMulti(ctx context.Context, keys []*datastore.Key, src interface{}) error {
	value := reflect.ValueOf(src)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		mc.Set(ctx, keys[i], item.Interface())
	}

	return nil
}

func (mc *cache) Contains(key *datastore.Key) bool {
	_, ok := mc.items[key.Encode()]
	return ok
}
