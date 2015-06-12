package horsefeather

import (
	"google.golang.org/appengine/datastore"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"
	"time"
)

//SaveStruct outputs a struct to an datastore.Property channel
func SaveStruct(src interface{}) ([]datastore.Property, error) {
	value := reflect.Indirect(reflect.ValueOf(src))
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		config := parseTag(sf.Tag.Get("hf"))

		switch config.Tag {
		case "timestamp":
			now := reflect.ValueOf(time.Now())
			value.FieldByName(sf.Name).Set(now)

		}
	}

	output, err := datastore.SaveStruct(src)
	if err != nil {
		return output, err
	}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		v := value.Field(i)
		config := parseTag(sf.Tag.Get("hf"))

		switch config.Tag {
		case "opaque":
			if data, err := json.Marshal(v.Interface()); err == nil {
				output = append(output, datastore.Property{Name: sf.Name, Value: data, NoIndex: true, Multiple: false})
				if config.CRC {
					hash := crc32.NewIEEE()
					fmt.Fprint(hash, data)
					value.FieldByName(sf.Name + "CRC").Set(reflect.ValueOf(fmt.Sprintf("%x", hash.Sum(nil))))
				}
			}
		case "compress":
			if data, err := json.Marshal(v.Interface()); err == nil {
				buf := &bytes.Buffer{}
				w := gzip.NewWriter(buf)
				w.Write(data)
				w.Close()
				output = append(output, datastore.Property{Name: sf.Name, Value: buf.Bytes(), NoIndex: true, Multiple: false})
				if config.CRC {
					hash := crc32.NewIEEE()
					fmt.Fprint(hash, buf.Bytes())
					value.FieldByName(sf.Name + "CRC").Set(reflect.ValueOf(fmt.Sprintf("%x", hash.Sum(nil))))
				}
			}
		}
	}

	return output, err
}
