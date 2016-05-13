package horsefeather

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/cogger.v1"
	"gopkg.in/cogger/cogger.v1/cogs"
	"gopkg.in/cogger/cogger.v1/order"
	"gopkg.in/cogger/cogger.v1/wait"
	"gopkg.in/cogger/stash.v1"
)

func GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	defer reset(ctx)
	value := reflect.Indirect(reflect.ValueOf(dst))
	result := newResult()
	workers := make([]cogger.Cog, 0, len(keys))
	createMemoryLoaders := cogs.Simple(ctx, func() error {
		if len(keys) != value.Len() {
			return ErrInvalidEntityType
		}
		for i := 0; i < len(keys); i++ {
			result.Set(keys[i], nil)
			func(i int) {

				found := false

				loadFromStash := cogs.Simple(ctx, func() error {
					encodedKey := keys[i].Encode()
					if stash.Has(ctx, encodedKey) {
						data := stash.Get(ctx, encodedKey)

						result.Set(keys[i], data)
						found = true
					}
					return nil
				})

				loadFromMC := order.If(ctx,
					func() bool { return !found && IsMemcacheAllowed(ctx) },
					cogs.Simple(ctx, func() error {
						item := value.Index(i).Type()
						if item.Kind() == reflect.Ptr {
							item = item.Elem()
						}
						data := reflect.New(item).Interface()
						err := mc(ctx).Get(ctx, keys[i], &data)
						if err == nil {
							result.Set(keys[i], data)
						}

						found = err == nil
						return nil
					}),
				)
				workers = append(workers, order.Series(ctx,
					loadFromStash,
					loadFromMC,
				))

			}(i)
		}
		return nil
	})

	executeMemoryLoaders := order.If(ctx,
		func() bool { return len(workers) > 0 },
		cogs.DeferredCreate(func() cogger.Cog {
			return order.Parallel(ctx, workers...)
		}),
	)

	var remainingKeys []*datastore.Key
	findMissingItems := cogs.Simple(ctx, func() error {
		remainingKeys = result.RemainingKeys()
		return nil
	})

	loadMissing := order.If(ctx,
		func() bool { return len(remainingKeys) > 0 && IsDatastoreAllowed(ctx) },
		cogs.Simple(ctx, func() error {
			l := len(remainingKeys)
			remainingItems := reflect.MakeSlice(value.Type(), l, l)

			err := ds(ctx).GetMulti(ctx, remainingKeys, remainingItems.Interface())
			for i := 0; i < remainingItems.Len(); i++ {
				result.Set(remainingKeys[i], remainingItems.Index(i).Interface())
			}

			return err
		}),
	)

	setResultsToDst := cogs.Simple(ctx, func() error {
		for i, key := range keys {
			item := result.Get(key)
			if item != nil && value.Index(i).Type().AssignableTo(reflect.TypeOf(item)) {
				value.Index(i).Set(reflect.ValueOf(item))
			}
		}
		return nil
	})

	errs := wait.Resolve(ctx,
		order.If(ctx,
			func() bool { return len(keys) > 0 },
			order.Series(ctx,
				createMemoryLoaders,
				executeMemoryLoaders,
				findMissingItems,
				loadMissing,
				setResultsToDst,
			),
		),
	)

	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}
