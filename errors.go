package horsefeather

import (
	"errors"
	"fmt"
)

var ErrIncompleteKey = errors.New("Key imcomplete but used a complete key.")

var ErrBadKey = errors.New("Key is corrupt")

var ErrNoContext = errors.New("context does not exist")

var ErrBadEntity = errors.New("could not save Entity")

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
