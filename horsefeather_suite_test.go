package horsefeather

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"

	"google.golang.org/appengine/datastore"
)

func TestHorsefeather(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Horsefeather Suite")
}

type simple struct {
	Name   string
	Age    int
	Date   time.Time
	Weight float64
}

type tagged struct {
	Name   string    `datastore:",noindex"`
	Age    int       `datastore:"age"`
	Date   time.Time `datastore:"-"`
	Weight float64   `datastore:"wt,noindex"`
}

type repeated struct {
	Name     string
	Children []string
}

type nested struct {
	Name    string
	Address addr
}

type addr struct {
	Street  string
	City    string
	State   string
	Zipcode string
}

type opaque struct {
	Name    string
	Address addr `datastore:"-" hf:"opaque"`
}

type compressed struct {
	Name    string
	Address addr `datastore:"-" hf:"compress"`
}

type opaqueslice struct {
	Name     string
	Children []string `datastore:"-" hf:"opaque"`
}

type opaquecrc struct {
	Name       string
	Address    *addr  `datastore:"-" hf:"opaque,crc"`
	AddressCRC string `datastore:",noindex"`
}

type compressedcrc struct {
	Name       string
	Address    *addr  `datastore:"-" hf:"compress,crc"`
	AddressCRC string `datastore:",noindex"`
}

type testCase struct {
	Name     string
	Object   interface{}
	Props    []datastore.Property
	Expected interface{}
}

var tests = []testCase{
	{
		Name: "simple",
		Object: &simple{
			Name:   "bobbi",
			Age:    25,
			Date:   time.Date(1998, 6, 12, 6, 17, 0, 0, time.Local),
			Weight: 130.5,
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Age", Value: int64(25), NoIndex: false, Multiple: false},
			datastore.Property{Name: "Date", Value: time.Date(1998, 6, 12, 6, 17, 0, 0, time.Local), NoIndex: false, Multiple: false},
			datastore.Property{Name: "Weight", Value: 130.5, NoIndex: false, Multiple: false},
		},
		Expected: &simple{
			Name:   "bobbi",
			Age:    25,
			Date:   time.Date(1998, 6, 12, 6, 17, 0, 0, time.Local),
			Weight: 130.5,
		},
	}, {
		Name: "with datastore tags",
		Object: &tagged{
			Name:   "bobbi",
			Age:    25,
			Date:   time.Date(1998, 6, 12, 6, 17, 0, 0, time.Local),
			Weight: 130.5,
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: true, Multiple: false},
			datastore.Property{Name: "age", Value: int64(25), NoIndex: false, Multiple: false},
			datastore.Property{Name: "wt", Value: 130.5, NoIndex: true, Multiple: false},
		},
		Expected: &tagged{
			Name:   "bobbi",
			Age:    25,
			Weight: 130.5,
		},
	}, {
		Name: "repeated",
		Object: &repeated{
			Name:     "bobbi",
			Children: []string{"ann", "joe"},
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Children", Value: "ann", NoIndex: false, Multiple: true},
			datastore.Property{Name: "Children", Value: "joe", NoIndex: false, Multiple: true},
		},
		Expected: &repeated{
			Name:     "bobbi",
			Children: []string{"ann", "joe"},
		},
	}, {
		Name: "nested",
		Object: &nested{
			Name: "bobbi",
			Address: addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Address.Street", Value: "1572 Sylvan Drive", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Address.City", Value: "York", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Address.State", Value: "PA", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Address.Zipcode", Value: "17402", NoIndex: false, Multiple: false},
		},
		Expected: &nested{
			Name: "bobbi",
			Address: addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
		},
	}, {
		Name: "opaque",
		Object: &opaque{
			Name: "bobbi",
			Address: addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Address", Value: []byte("{\"Street\":\"1572 Sylvan Drive\",\"City\":\"York\",\"State\":\"PA\",\"Zipcode\":\"17402\"}"), NoIndex: true, Multiple: false},
		},
		Expected: &opaque{
			Name: "bobbi",
			Address: addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
		},
	}, {
		Name: "compressed",
		Object: &compressed{
			Name: "bobbi",
			Address: addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Address", Value: []byte("\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xaaV\n.)JM-Q\xb2R2457R\b\xae\xcc)K\xccSp)\xca,KU\xd2Qr\xce,\xa9\x04JE\xe6\x17e\x03y\xc1%\x89%\xa9@n\x80#\x90\x13\x95Y\x90\x9c\x9f\x02\xe2\x1a\x9a\x9b\x18\x18)\xd5\x02\x02\x00\x00\xff\xff\x95\xbeɼK\x00\x00\x00"), NoIndex: true, Multiple: false},
		},
		Expected: &compressed{
			Name: "bobbi",
			Address: addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
		},
	}, {
		Name: "opaqueslice",
		Object: &opaqueslice{
			Name:     "bobbi",
			Children: []string{"bob", "ann"},
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "Children", Value: []byte("[\"bob\",\"ann\"]"), NoIndex: true, Multiple: false},
		},
		Expected: &opaqueslice{
			Name:     "bobbi",
			Children: []string{"bob", "ann"},
		},
	}, {
		Name: "opaquecrc",
		Object: &opaquecrc{
			Name: "bobbi",
			Address: &addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
			AddressCRC: "b08e4eed",
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "AddressCRC", Value: "b08e4eed", NoIndex: true, Multiple: false},
			datastore.Property{Name: "Address", Value: []byte("{\"Street\":\"1572 Sylvan Drive\",\"City\":\"York\",\"State\":\"PA\",\"Zipcode\":\"17402\"}"), NoIndex: true, Multiple: false},
		},
		Expected: &opaquecrc{
			Name: "bobbi",
			Address: &addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
			AddressCRC: "b08e4eed",
		},
	}, {
		Name: "compressedcrc",
		Object: &compressedcrc{
			Name: "bobbi",
			Address: &addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
			AddressCRC: "0d109e7e",
		},
		Props: []datastore.Property{
			datastore.Property{Name: "Name", Value: "bobbi", NoIndex: false, Multiple: false},
			datastore.Property{Name: "AddressCRC", Value: "0d109e7e", NoIndex: true, Multiple: false},
			datastore.Property{Name: "Address", Value: []byte("\x1f\x8b\b\x00\x00\tn\x88\x00\xff\xaaV\n.)JM-Q\xb2R2457R\b\xae\xcc)K\xccSp)\xca,KU\xd2Qr\xce,\xa9\x04JE\xe6\x17e\x03y\xc1%\x89%\xa9@n\x80#\x90\x13\x95Y\x90\x9c\x9f\x02\xe2\x1a\x9a\x9b\x18\x18)\xd5\x02\x02\x00\x00\xff\xff\x95\xbeɼK\x00\x00\x00"), NoIndex: true, Multiple: false},
		},
		Expected: &compressedcrc{
			Name: "bobbi",
			Address: &addr{
				Street:  "1572 Sylvan Drive",
				City:    "York",
				State:   "PA",
				Zipcode: "17402",
			},
			AddressCRC: "0d109e7e",
		},
	},
}
