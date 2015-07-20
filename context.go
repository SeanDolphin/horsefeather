package horsefeather

import "golang.org/x/net/context"

var contextKey = "horsefeatherkey"

func c(ctx context.Context) *container {
	value := ctx.Value(&contextKey)
	if value == nil {
		return &container{}
	}

	box := value.(*container)

	return box
}

func setC(ctx context.Context, box *container) context.Context {
	ctx = context.WithValue(ctx, &contextKey, box)
	return ctx
}

type container struct {
	mc   Memcache
	noMC bool

	ds   Datastore
	noDS bool
}

func reset(ctx context.Context) {
	box := c(ctx)
	box.noDS = false
	box.noMC = false
	setC(ctx, box)
}
