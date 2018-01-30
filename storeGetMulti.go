package horsefeather

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func GetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	defer reset(ctx)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	value := reflect.Indirect(reflect.ValueOf(dst))
	result := newResult()

	var errs = []error{}
	if len(keys) != value.Len() {
		return ErrInvalidEntityType
	}

	for i := 0; i < len(keys); i++ {
		result.Set(ctx, keys[i], nil)

		if IsMemcacheAllowed(ctx) {
			item := value.Index(i).Type()
			if item.Kind() == reflect.Ptr {
				item = item.Elem()
			}
			data := reflect.New(item).Interface()
			if err := mc(ctx).Get(ctx, keys[i], &data); err == nil {
				result.Set(ctx, keys[i], data)
			}
		}
	}

	if IsDatastoreAllowed(ctx) {
		var remainingKeys = result.RemainingKeys()
		l := len(remainingKeys)

		if l > 0 {
			remainingItems := reflect.MakeSlice(value.Type(), l, l)

			if err := ds(ctx).GetMulti(ctx, remainingKeys, remainingItems.Interface()); err == nil {
				for i := 0; i < l; i++ {
					result.Set(ctx, remainingKeys[i], remainingItems.Index(i).Interface())
				}
			} else {
				errs = append(errs, err)
			}

		}
	}

	for i, key := range keys {
		item := result.Get(ctx, key)
		if item != nil && value.Index(i).Type().AssignableTo(reflect.TypeOf(item)) {
			value.Index(i).Set(reflect.ValueOf(item))
		}
	}

	if len(errs) > 0 {
		return ErrMulti(errs)
	}
	return nil
}
