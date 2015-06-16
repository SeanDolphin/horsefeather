package horsefeather

import "golang.org/x/net/context"

var contextKey = "horsefeatherkey"

func c(ctx context.Context) *container {
	value := ctx.Value(&contextKey)
	if value == nil {
		return &container{}
	}

	box, _ := value.(*container)
	return box
}

func setC(ctx context.Context, box *container) context.Context {
	return context.WithValue(ctx, &contextKey, box)
}

type container struct {
	mc Memcache
	ds Datastore
}
