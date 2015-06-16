package horsefeather

import (
	"errors"
	"fmt"
)

var ErrIncompleteKey = errors.New("Key imcomplete but used a complete key.")

var ErrNoContext = errors.New("context does not exist")

// ErrInvalidEntityType is returned when functions like Get or Next are
// passed a dst or src argument of invalid type.
var ErrInvalidEntityType = errors.New("datastore: invalid entity type")

// ErrInvalidKey is returned when an invalid key is presented.
var ErrInvalidKey = errors.New("datastore: invalid key")

// ErrNoSuchEntity is returned when no entity was found for a given key.
var ErrNoSuchEntity = errors.New("datastore: no such entity")

type ErrMulti []error

func (m ErrMulti) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}
