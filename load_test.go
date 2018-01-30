package horsefeather

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadStruct", func() {
	for _, test := range tests {
		runLoadTest(test)
	}
})

func runLoadTest(test testCase) {
	It(fmt.Sprintf("should load %s structs", test.Name), func() {
		dst := reflect.New(reflect.TypeOf(test.Object).Elem()).Interface()
		err := LoadStruct(dst, test.Props)
		Expect(err).ToNot(HaveOccurred())
		Expect(dst).To(BeEquivalentTo(test.Expected))
	})
}
