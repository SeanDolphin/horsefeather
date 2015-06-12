package horsefeather

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func mc(ctx context.Context) *mem {
	return &mem{
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
	}
}

type mem struct {
	Codec memcache.Codec
}

func (mc *mem) Delete(ctx context.Context, key *datastore.Key) error {
	return memcache.Delete(ctx, key.Encode())
}

func (mc *mem) DeleteMulti(ctx context.Context, keys []*datastore.Key) error {
	simpleKeys := make([]string, len(keys))
	for i, key := range keys {
		simpleKeys[i] = key.Encode()
	}

	return memcache.DeleteMulti(ctx, simpleKeys)
}

func (mc *mem) Get(ctx context.Context, key *datastore.Key, dst interface{}) error {
	_, err := mc.Codec.Get(ctx, key.Encode(), dst)
	return err
}

func (mc *mem) Set(ctx context.Context, key *datastore.Key, dst interface{}) error {
	return mc.Codec.Set(ctx, &memcache.Item{
		Key:    key.Encode(),
		Object: dst,
	})
}

func (mc *mem) SetMulti(ctx context.Context, keys []*datastore.Key, dst interface{}) error {
	return nil
}
