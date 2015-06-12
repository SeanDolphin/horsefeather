package horsefeather

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrMulti", func() {
	It("should return \"(0 errors)\" when there are no errors", func() {
		errs := ErrMulti{}
		Expect(errs.Error()).To(Equal("(0 errors)"))
	})

	It("should return the error when there is only 1 error", func() {
		errMsg := "test message"
		errs := ErrMulti{errors.New(errMsg)}
		Expect(errs.Error()).To(Equal(errMsg))
	})

	It("should return the first error and \"(and 1 other error)\" when there are 2 errors.", func() {
		errMsg := "test message"
		errs := ErrMulti{errors.New(errMsg), errors.New(errMsg)}
		Expect(errs.Error()).To(Equal(fmt.Sprintf("%s (and 1 other error)", errMsg)))
	})

	It("should return the first error and \"(and # other errors)\" where are are more then 2 errors.", func() {
		errMsg := "test message"
		errs := ErrMulti{errors.New(errMsg), errors.New(errMsg)}
		for i := 1; i < 10; i++ {
			errs = append(errs, errors.New(errMsg))
			Expect(errs.Error()).To(Equal(fmt.Sprintf("%s (and %d other errors)", errMsg, i+1)))
		}

	})
})
