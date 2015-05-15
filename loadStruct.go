package horsefeather

import (
	"google.golang.org/appengine/datastore"

	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"reflect"
)

//LoadStruct helps load a struct from a datastore.Property channel
func LoadStruct(dst interface{}, input []datastore.Property) error {
	filtered := []datastore.Property{}

	value := reflect.Indirect(reflect.ValueOf(dst))
	if value.IsValid() {
		tValue := value.Type()
		for _, prop := range input {
			if !prop.Multiple {

				sf, ok := tValue.FieldByName(prop.Name)
				config := tag{}

				if ok {
					config = parseTag(sf.Tag.Get("hf"))
				}

				switch config.Tag {
				case "opaque", "compress":
					t := sf.Type

					var v reflect.Value
					switch t.Kind() {
					case reflect.Ptr:
						v = reflect.New(t.Elem())
					case reflect.Slice:
						v = reflect.New(reflect.SliceOf(t.Elem()))
					default:
						v = reflect.New(t)
					}
					data := prop.Value.([]byte)
					if config.Tag == "compress" {
						first := data[0]
						last := data[len(data)-1]
						if !(first == 123 && last == 125) && !(first == 91 && last == 93) {
							buf := bytes.NewBuffer(data)
							output := &bytes.Buffer{}
							r, _ := gzip.NewReader(buf)
							defer r.Close()
							io.Copy(output, r)
							data = output.Bytes()
						}
					}
					if err := json.Unmarshal(data, v.Interface()); err == nil {
						if v.Kind() == reflect.Ptr && sf.Type.Kind() != reflect.Ptr {
							v = v.Elem()
						}
						value.FieldByName(prop.Name).Set(v)

					}
				default:
					filtered = append(filtered, prop)
				}
			} else {
				filtered = append(filtered, prop)
			}
		}
	}
	return datastore.LoadStruct(dst, filtered)
}
