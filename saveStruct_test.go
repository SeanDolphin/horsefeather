package horsefeather

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SaveStruct", func() {
	for _, test := range tests {
		runSaveTest(test)
	}

	It("should error when an item can not be saved to the datastore", func() {
		type nest struct {
			Name string
			Nest *nest
		}

		_, err := SaveStruct(nest{
			Name: "test1",
			Nest: &nest{
				Name: "test2",
				Nest: &nest{
					Name: "test3",
					Nest: nil,
				},
			},
		})
		Expect(err).To(HaveOccurred())
	})

	It("should timestamp structs when they are tagged and saved", func() {
		type timestamptest struct {
			Name string
			Time time.Time `hf:"timestamp"`
		}

		tst := timestamptest{Name: "test"}
		_, err := SaveStruct(&tst)
		Expect(err).ToNot(HaveOccurred())
		Expect(time.Now().Sub(tst.Time).Seconds() < 1).To(BeTrue())
	})
})

func runSaveTest(test testCase) {
	It(fmt.Sprintf("should save %s structs", test.Name), func() {

		output, err := SaveStruct(test.Object)

		Expect(err).ToNot(HaveOccurred())
		Expect(output).To(HaveLen(len(test.Props)))
		Expect(output).To(BeEquivalentTo(test.Props))
	})
}
