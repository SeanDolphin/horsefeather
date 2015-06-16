package test

import (
	"encoding/json"
	"errors"
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func NewStore() *store {
	return &store{
		items: map[string][]byte{},
	}
}

type store struct {
	items map[string][]byte
}

func (ds *store) Delete(ctx context.Context, key *datastore.Key) error {
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
	data, ok := ds.items[key.Encode()]
	if ok {
		return json.Unmarshal(data, &dst)
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
		value.Index(i).Set(reflect.ValueOf(item))
	}
	return nil
}

func (ds *store) Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	if src == nil {
		return key, errors.New("item is nil")
	}
	data, err := json.Marshal(src)
	if err == nil {
		ds.items[key.Encode()] = data
	}
	return key, err
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
	_, ok := ds.items[key.Encode()]
	return ok
}

func (ds *store) Clear() {
	ds.items = map[string][]byte{}
}
