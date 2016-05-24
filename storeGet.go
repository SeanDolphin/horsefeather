package horsefeather

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/cogger.v1/cogs"
	"gopkg.in/cogger/cogger.v1/order"
	"gopkg.in/cogger/cogger.v1/wait"
	"gopkg.in/cogger/stash.v1"
)

func Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	defer reset(ctx)
	found := false
	errs := wait.Resolve(ctx,
		order.Series(ctx,
			cogs.Simple(ctx, func() error {
				encodedKey := key.Encode()
				if stash.Has(ctx, encodedKey) {
					found = true
					item := stash.Get(ctx, encodedKey)
					reflect.Indirect(reflect.ValueOf(dst)).Set(reflect.ValueOf(item))
				}
				return nil
			}),
			order.If(ctx,
				func() bool { return !found && IsMemcacheAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					err := mc(ctx).Get(ctx, key, dst)
					found = err == nil
					return nil
				}),
			),
			order.If(ctx,
				func() bool { return !found && IsDatastoreAllowed(ctx) },
				cogs.Simple(ctx, func() error {
					err := ds(ctx).Get(ctx, key, dst)
					if err == nil {
						mc(ctx).Set(ctx, key, dst)
					}
					return err
				}),
			),
		),
	)

	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}
