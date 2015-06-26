package test

import (
	"errors"
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func NewCache() *cache {
	return &cache{
		items: map[string]interface{}{},
	}
}

type cache struct {
	items map[string]interface{}
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
	item, ok := mc.items[key.Encode()]
	if ok {
		reflect.Indirect(reflect.ValueOf(dst)).Set(reflect.ValueOf(item))
		return nil
	}
	return errors.New("entity does not exist")
}

func (mc *cache) Set(ctx context.Context, key *datastore.Key, src interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(src))
	if !value.IsValid() {
		return errors.New("invalid entity")
	}
	mc.items[key.Encode()] = value.Interface()
	return nil
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

func (mc *cache) Clear() {
	mc.items = map[string]interface{}{}
}
