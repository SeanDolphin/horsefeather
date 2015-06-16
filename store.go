package horsefeather

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/cogger.v1"
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

func Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	found := false
	errs := wait.Resolve(ctx,
		order.Series(ctx,
			cogs.Simple(ctx, func() error {
				err := mc(ctx).Get(ctx, key, dst)
				found = err == nil
				return nil
			}),
			order.If(ctx,
				func() bool { return !found },
				cogs.Simple(ctx, func() error {
					return ds(ctx).Get(ctx, key, dst)
				}),
			),
		),
	)

	if len(errs) > 0 {
		return ErrMulti(errs)
	}
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

func GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	value := reflect.Indirect(reflect.ValueOf(dst))

	result := map[*datastore.Key]interface{}{}
	workers := make([]cogger.Cog, 0, len(keys))
	createMemoryLoaders := cogs.Simple(ctx, func() error {
		for i := 0; i < len(keys); i++ {
			result[keys[i]] = nil
			func(i int) {
				item := value.Index(i)
				workers = append(workers, cogs.Simple(ctx, func() error {
					data := item.Interface()
					err := mc(ctx).Get(ctx, keys[i], &data)

					if err == nil {
						result[keys[i]] = data
					}
					return nil
				}))
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

	remainingKeys := make([]*datastore.Key, 0, len(keys))
	findMissingItems := cogs.Simple(ctx, func() error {
		for key, item := range result {
			if item == nil {
				remainingKeys = append(remainingKeys, key)
			}
		}

		return nil
	})

	loadMissing := order.If(ctx,
		func() bool { return len(remainingKeys) > 0 },
		cogs.Simple(ctx, func() error {
			l := len(remainingKeys)
			remainingItems := reflect.MakeSlice(value.Type(), l, l)

			err := ds(ctx).GetMulti(ctx, remainingKeys, remainingItems.Interface())
			for i := 0; i < remainingItems.Len(); i++ {
				result[remainingKeys[i]] = remainingItems.Index(i).Interface()
			}

			return err
		}),
	)

	setResultsToDst := cogs.Simple(ctx, func() error {
		for i, key := range keys {
			item := result[key]
			if item != nil {
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
