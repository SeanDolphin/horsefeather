package hfae

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/pquerna/ffjson/ffjson"

	"google.golang.org/appengine/memcache"
	"gopkg.in/SeanDolphin/horsefeather.v1"
)

const onemb = 1000000

var Codec = memcache.Codec{
	Marshal: func(src interface{}) ([]byte, error) {
		data, err := ffjson.Marshal(src)
		if err != nil {
			return data, err
		}
		buf := &bytes.Buffer{}
		w := gzip.NewWriter(buf)
		w.Write(data)
		w.Close()

		ffjson.Pool(data)
		if buf.Len() > onemb {
			return []byte{}, horsefeather.ErrEntityToLarge
		}
		return buf.Bytes(), err
	},
	Unmarshal: func(data []byte, dst interface{}) error {
		buf := bytes.NewBuffer(data)
		output := &bytes.Buffer{}
		output.Grow(buf.Len())
		r, _ := gzip.NewReader(buf)
		defer r.Close()
		io.Copy(output, r)

		ffjson.Pool(data)
		return ffjson.Unmarshal(output.Bytes(), dst)
	},
}
