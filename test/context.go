package test

import (
	"net/http"

	"golang.org/x/net/context"
	"gopkg.in/SeanDolphin/horsefeather.v1"
)

func Add(ctx context.Context, req *http.Request) context.Context {
	ctx = Set(ctx)
	return ctx
}

func Set(ctx context.Context) context.Context {
	ctx = horsefeather.AddMemcache(ctx, NewCache())
	ctx = horsefeather.AddDatastore(ctx, NewStore())
	return ctx
}
