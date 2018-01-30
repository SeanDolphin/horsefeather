package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	defer reset(ctx)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var keyResult *datastore.Key
	var errs = []error{}

	if IsMemcacheAllowed(ctx) {
		if err := mc(ctx).Set(ctx, key, src); err != nil {
			if !IsDatastoreAllowed(ctx) {
				errs = append(errs, err)
			}
		}
	}

	if IsDatastoreAllowed(ctx) {
		createdKey, err := ds(ctx).Put(ctx, key, src)
		if err == nil {
			keyResult = createdKey
		} else {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return keyResult, ErrMulti(errs)
	}
	return keyResult, nil
}

func PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	defer reset(ctx)
	var keyResult []*datastore.Key
	var errs = []error{}

	if IsDatastoreAllowed(ctx) {
		createdKey, err := ds(ctx).PutMulti(ctx, keys, src)
		if err == nil {
			keyResult = createdKey
		} else {
			errs = append(errs, err)
		}
	}
	if IsMemcacheAllowed(ctx) {
		if err := mc(ctx).SetMulti(ctx, keys, src); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return keyResult, ErrMulti(errs)
	}
	return keyResult, nil
}
