package test

import (
	"errors"
	"reflect"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func NewStore() *store {
	return &store{
		items: map[string]interface{}{},
	}
}

type store struct {
	sync.RWMutex
	items map[string][]byte
}

func (ds *store) Len() int {
	ds.RLock()
	defer ds.RUnlock()
	return len(ds.items)
}

func (ds *store) Delete(ctx context.Context, key *datastore.Key) error {
	ds.Lock()
	defer ds.Unlock()

	if _, ok := ds.items[key.Encode()]; !ok {
		return errors.New("does not exist")
	}
	delete(ds.items, key.Encode())
	return nil
}

func (ds *store) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	for _, key := range keys {
		err := ds.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *store) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	ds.RLock()
	defer ds.RUnlock()

	data, ok := ds.items[key.Encode()]
	if ok {
		reflect.Indirect(reflect.ValueOf(dst)).Set(reflect.Indirect(reflect.ValueOf(data)))
		// reflect.ValueOf(dst).Set(reflect.Indirect(reflect.ValueOf(data)))
		return nil
	}

	return errors.New("entity does not exist")
}

func (ds *store) GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	value := reflect.ValueOf(dst)
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		err := ds.Get(ctx, keys[i], &item)
		if err != nil {
			return err
		}
		valueOf := reflect.ValueOf(item)
		v := value.Index(i)
		if v.Kind() == reflect.Ptr {
			v.Set(reflect.New(v.Type().Elem()))
		} else {
			v.Set(valueOf)
		}

	}
	return nil
}

func (ds *store) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	ds.Lock()
	defer ds.Unlock()
	if src == nil {
		return key, errors.New("item is nil")
	}

	ds.items[key.Encode()] = src

	return key, nil
}

func (ds *store) PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	value := reflect.ValueOf(src)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Len() != len(keys) {
		return keys, errors.New("length and keys do not match")
	}

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		_, err := ds.Put(ctx, keys[i], item.Interface())
		if err != nil {
			return keys, err
		}
	}

	return keys, nil
}

func (ds *store) Contains(key *datastore.Key) bool {
	ds.RLock()
	defer ds.RUnlock()
	_, ok := ds.items[key.Encode()]
	return ok
}

func (ds *store) Clear() {
	ds.Lock()
	defer ds.Unlock()
	ds.items = map[string][]byte{}
}
