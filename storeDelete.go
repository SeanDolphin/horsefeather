package horsefeather

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/cogger.v1/cogs"
	"gopkg.in/cogger/cogger.v1/order"
	"gopkg.in/cogger/cogger.v1/wait"
	"gopkg.in/cogger/stash.v1"
)

func Delete(ctx context.Context, key *datastore.Key) error {
	defer reset(ctx)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			cogs.Simple(ctx, func() error {
				stash.Remove(ctx, key.Encode())
				return nil
			}),
			order.If(ctx,
				func() bool { return IsDatastoreAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					return ds(ctx).Delete(ctx, key)
				}),
			),
			order.If(ctx,
				func() bool { return IsMemcacheAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					return mc(ctx).Delete(ctx, key)
				}),
			),
		),
	)
	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}

func DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	defer reset(ctx)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	errs := wait.Resolve(ctx,
		order.Parallel(ctx,
			cogs.Simple(ctx, func() error {
				for _, key := range keys {
					stash.Remove(ctx, key.Encode())
				}
				return nil
			}),
			order.If(ctx,
				func() bool { return IsDatastoreAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					return ds(ctx).DeleteMulti(ctx, keys)
				}),
			),
			order.If(ctx,
				func() bool { return IsMemcacheAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					return mc(ctx).DeleteMulti(ctx, keys)
				}),
			),
		),
	)
	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}
