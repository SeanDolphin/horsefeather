package horsefeather

import (
	"reflect"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"gopkg.in/cogger/stash.v1"
)

func Prefetch(ctx context.Context, keys []*datastore.Key, dst interface{}) context.Context {
	defer reset(ctx)
	if len(keys) > 0 {
		err := GetMulti(ctx, keys, dst)
		if err == nil {
			value := reflect.Indirect(reflect.ValueOf(dst))
			for i := 0; i < len(keys); i++ {
				item := value.Index(i)
				ctx = stash.Set(ctx, keys[i].Encode(), item.Interface())
			}
		}
	}
	return ctx
}
