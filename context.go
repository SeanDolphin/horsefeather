package horsefeather

import (
	"sync"

	"golang.org/x/net/context"
)

var contextKey = "horsefeatherkey"

func c(ctx context.Context) (box *container) {
	value := ctx.Value(&contextKey)
	if value == nil {
		box = &container{}
		setC(ctx, box)
	} else {
		box = value.(*container)
	}

	return box
}

func setC(ctx context.Context, box *container) context.Context {
	ctx = context.WithValue(ctx, &contextKey, box)
	return ctx
}

type container struct {
	sync.RWMutex

	mc   Memcache
	noMC bool

	ds   Datastore
	noDS bool
}

func reset(ctx context.Context) {
	box := c(ctx)
	box.Lock()
	defer box.Unlock()

	box.noDS = false
	box.noMC = false
	setC(ctx, box)
}
