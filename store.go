package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/cogger.v1/cogs"
	"gopkg.in/cogger/cogger.v1/order"
	"gopkg.in/cogger/cogger.v1/wait"
)

func Delete(ctx context.Context, key *datastore.Key) error {
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			cogs.Simple(ctx, func() error {
				return ds(ctx).Delete(ctx, key)
			}),
			cogs.Simple(ctx, func() error {
				return mc(ctx).Delete(ctx, key)
			}),
		),
	)
	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}

func DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			cogs.Simple(ctx, func() error {
				return ds(ctx).DeleteMulti(ctx, keys)
			}),
			cogs.Simple(ctx, func() error {
				return mc(ctx).DeleteMulti(ctx, keys)
			}),
		),
	)
	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}

func Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	return <-wait.Any(ctx,
		cogs.Simple(ctx, func() error {
			return mc(ctx).Get(ctx, key, dst)
		}),
		cogs.Simple(ctx, func() error {
			return ds(ctx).Get(ctx, key, dst)
		}),
	).Do(ctx)
}

func GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	return nil
}

func Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	var keyResult *datastore.Key
	var err error
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			cogs.Simple(ctx, func() error {
				keyResult, err = ds(ctx).Put(ctx, key, src)
				return err
			}),
			cogs.Simple(ctx, func() error {
				return mc(ctx).Set(ctx, key, src)
			}),
		),
	)
	if len(errs) > 0 {
		return keyResult, ErrMulti(errs)
	}
	return keyResult, nil
}

func PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	var keyResult []*datastore.Key
	var err error
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			cogs.Simple(ctx, func() error {
				keyResult, err = ds(ctx).PutMulti(ctx, keys, src)
				return err
			}),
			cogs.Simple(ctx, func() error {
				return mc(ctx).SetMulti(ctx, keys, src)
			}),
		),
	)
	if len(errs) > 0 {
		return keyResult, ErrMulti(errs)
	}
	return keyResult, nil
}

func Set(ctx context.Context, key *datastore.Key, src interface{}) error {
	return mc(ctx).Set(ctx, key, src)
}

func SetMulti(ctx context.Context, keys []*datastore.Key, src interface{}) error {
	return mc(ctx).SetMulti(ctx, keys, src)
}
