# horsefeather

[![GoDoc](https://godoc.org/github.com/SeanDolphin/horsefeather?status.png)](http://godoc.org/github.com/SeanDolphin/horsefeather)  
[![Build Status](https://travis-ci.org/SeanDolphin/horsefeather.svg?branch=master)](https://travis-ci.org/SeanDolphin/horsefeather)  
[![Coverage Status](https://coveralls.io/repos/SeanDolphin/horsefeather/badge.svg)](https://coveralls.io/r/SeanDolphin/horsefeather)  
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

horsefeather is a [Google Appengine Engine](https://github.com/golang/appengine) extension that adds functionality to though tags to handle saving and loading datastore models differently.

**Usage**

~~~ go
package somepackage

import (
	"github.com/SeanDolphin/horsefeather"
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
