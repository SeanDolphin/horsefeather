# horsefeather

[![GoDoc](https://godoc.org/github.com/SeanDolphin/horsefeather?status.png)](http://godoc.org/github.com/SeanDolphin/horsefeather)  
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

horsefeather is a [Google Appengine Engine](https://github.com/golang/appengine) extension that adds functionality through tags to handle saving and loading datastore models differently.  horsefeather as of v1.1 contains helper functions for basic manipulation of datastore operations.  It operates like NDB from Python Appengine.

## Installation

The import path for the package is *gopkg.in/SeanDolphin/horsefeather.v1*.

To install it, run:

    go get gopkg.in/SeanDolphin/horsefeather.v1

## Usage

### Tags
~~~ go
package somepackage

import (
	"gopkg.in/SeanDolphin/horsefeather.v1"
	"google.golang.org/appengine/datastore"

	"time"
)

type Person struct{
	Key    		*datastore.Key
	Name 		FullName 		`datastore:"-" hf:"opaque"`			//saves this property as a json string
	Address    	*Address  		`datastore:"-" hf:"compress,crc"`	//compresses this property and calculates and crc
	AddressCRC	string			`datastore:",noindex"`				//set by horsefeather crc above
	Updated 	time.TIme 		`hf:"timestamp"`					//updated on save with current timestamp
}

func (person *Content) Load(input []datastore.Property) error {
	return horsefeather.LoadStruct(person, input)
}

func (person *Content) Save() ([]datastore.Property, error) {
	return horsefeather.SaveStruct(person)
}

type FullName struct{
	First	string
	Last 	string
}

type Address struct {
	Street  string
	City    string
	State   string
	Zipcode string
}
~~~

### Helpers

Setting up is as simple as calling hfae on the context you are using.  

~~~ go
	ctx = hfae.Set(ctx) 
~~~

The calls then work as normal calls would work except it works with memcache, datastore and a local cache in unison.
~~~ go
package somepackage

type Person struct{
	Key    		*datastore.Key
	Name 		FullName 						
	Updated 	time.TIme 		
}

func DoSometask(ctx context.Context){
	var person Person
	err := hs.Get(ctx, key, &person)
}
~~~
