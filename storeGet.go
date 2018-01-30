package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	defer reset(ctx)
	if IsMemcacheAllowed(ctx) {
		if err := mc(ctx).Get(ctx, key, dst); err == nil {
			return nil
		}
	}

	if IsDatastoreAllowed(ctx) {
		return ds(ctx).Get(ctx, key, dst)
	}

	return ErrNoSuchEntity
}
