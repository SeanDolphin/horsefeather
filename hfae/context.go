package hfae

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
	"gopkg.in/SeanDolphin/horsefeather.v1"
)

func Add(ctx context.Context, req *http.Request) context.Context {
	ctx = Set(ctx)
	return ctx
}

func Set(ctx context.Context) context.Context {
	ctx = horsefeather.AddMemcache(ctx, &cache{
		Codec: memcache.Codec{
			Marshal: func(src interface{}) ([]byte, error) {
				data, err := json.Marshal(src)
				if err != nil {
					return data, err
				}
				buf := &bytes.Buffer{}
				w := gzip.NewWriter(buf)
				w.Write(data)
				w.Close()
				return buf.Bytes(), err
			},
			Unmarshal: func(data []byte, dst interface{}) error {
				buf := bytes.NewBuffer(data)
				output := &bytes.Buffer{}
				r, _ := gzip.NewReader(buf)
				defer r.Close()
				io.Copy(output, r)
				data = output.Bytes()
				return json.Unmarshal(data, dst)
			},
		},
	})
	ctx = horsefeather.AddDatastore(ctx, &store{})

	return ctx
}
