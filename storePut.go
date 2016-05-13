package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/cogger.v1/cogs"
	"gopkg.in/cogger/cogger.v1/order"
	"gopkg.in/cogger/cogger.v1/wait"
)

func Put(ctx context.Context, key *datastore.Key, src interface{}) (*datastore.Key, error) {
	defer reset(ctx)
	var keyResult *datastore.Key
	var err error
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			order.If(ctx,
				func() bool { return IsDatastoreAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					keyResult, err = ds(ctx).Put(ctx, key, src)
					return err
				}),
			),
			order.If(ctx,
				func() bool { return IsMemcacheAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					return mc(ctx).Set(ctx, key, src)
				}),
			),
		),
	)

	if len(errs) > 0 {
		return keyResult, ErrMulti(errs)
	}
	return keyResult, nil
}

func PutMulti(ctx context.Context, keys []*datastore.Key, src interface{}) ([]*datastore.Key, error) {
	defer reset(ctx)
	var keyResult []*datastore.Key
	var err error
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			order.If(ctx,
				func() bool { return IsDatastoreAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					keyResult, err = ds(ctx).PutMulti(ctx, keys, src)
					return err
				}),
			),
			order.If(ctx,
				func() bool { return IsMemcacheAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					return mc(ctx).SetMulti(ctx, keys, src)
				}),
			),
		),
	)
	if len(errs) > 0 {
		return keyResult, ErrMulti(errs)
	}
	return keyResult, nil
}
