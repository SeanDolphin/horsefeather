package horsefeather

import (
	"github.com/SeanDolphin/horsefeather/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var encodedKeys = []string{
	"aghnbGlicmFyeXIMCxIGUGVyc29uGAEM",
	"aghnbGlicmFyeXIdCxIFR3JhcGgiEmdyYXBoOjctZGF5LWFjdGl2ZQw",
	"aghnbGlicmFyeXIhCxIJV29yZEluZGV4GIChPgwLEglXb3JkSW5kZXgYiQgM",
}

var _ = Describe("Store", func() {
	var ctx context.Context
	var cache = test.NewCache()
	var store = test.NewStore()

	var keys []*datastore.Key
	for _, encodedKey := range encodedKeys {
		key, _ := datastore.DecodeKey(encodedKey)
		keys = append(keys, key)
	}

	It("should panic when memcache is not present", func() {
		Expect(func() { mc(context.Background()) }).To(Panic())
	})

	It("should panic when datastore is not present", func() {
		Expect(func() { ds(context.Background()) }).To(Panic())
	})

	Context("when caching context are present", func() {
		BeforeEach(func() {
			ctx = AddMemcache(context.Background(), cache)
			ctx = AddDatastore(ctx, store)
		})

		AfterEach(func() {
			store.Clear()
			cache.Clear()
		})

		Context("when dealing with single keys", func() {
			var data = "bob"
			Context("when deleteing", func() {
				It("should delete the key", func() {
					for _, key := range keys {
						_, err := Put(ctx, key, &data)
						Expect(err).ToNot(HaveOccurred())
						Expect(cache.Contains(key)).To(BeTrue(), "memcache")
						Expect(store.Contains(key)).To(BeTrue(), "datastore")

						err = Delete(ctx, key)
						Expect(err).ToNot(HaveOccurred())
						Expect(cache.Contains(key)).To(BeFalse(), "memcache")
						Expect(store.Contains(key)).To(BeFalse(), "datastore")
					}
				})

				It("should error when there is nothing to delete", func() {
					for _, key := range keys {
						Expect(Delete(ctx, key)).To(HaveOccurred())
					}
				})
			})

			Context("when getting", func() {
				It("should get the key", func() {
					for _, key := range keys {
						_, err := Put(ctx, key, &data)
						Expect(err).ToNot(HaveOccurred())

						var result string
						err = Get(ctx, key, &result)
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal(data))
					}
				})

				It("should get the key when it is not present in memcache", func() {
					for _, key := range keys {
						_, err := Put(ctx, key, &data)
						Expect(err).ToNot(HaveOccurred())
						cache.Delete(ctx, key)

						var result string
						err = Get(ctx, key, &result)
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal(data))
					}
				})

				It("should get the key when it is not present in the datastore", func() {
					for _, key := range keys {
						_, err := Put(ctx, key, &data)
						Expect(err).ToNot(HaveOccurred())
						store.Delete(ctx, key)

						var result string
						err = Get(ctx, key, &result)
						Expect(err).ToNot(HaveOccurred())
						Expect(result).To(Equal(data))
					}
				})

				It("should error when trying to get something that does not exist", func() {
					for _, key := range keys {
						var result string
						Expect(Get(ctx, key, &result)).To(HaveOccurred())
					}
				})
			})

			Context("whenning putting a key", func() {
				It("should put the data", func() {
					for _, key := range keys {
						_, err := Put(ctx, key, &data)
						Expect(err).ToNot(HaveOccurred())
						Expect(cache.Contains(key)).To(BeTrue(), "memcache")
						Expect(store.Contains(key)).To(BeTrue(), "datastore")
					}
				})

				It("should error when trying to save nil", func() {
					for _, key := range keys {
						_, err := Put(ctx, key, nil)
						Expect(err).To(HaveOccurred())
					}
				})
			})

		})

		Context("when dealing with multiple keys", func() {
			var data = []string{"t1", "t2", "t3"}
			Context("when deleting", func() {
				It("should delete all the keys", func() {
					_, err := PutMulti(ctx, keys, data)
					Expect(err).ToNot(HaveOccurred())
					for _, key := range keys {
						Expect(cache.Contains(key)).To(BeTrue(), "memcache")
						Expect(store.Contains(key)).To(BeTrue(), "store")
					}

					err = DeleteMulti(ctx, keys)
					Expect(err).ToNot(HaveOccurred())
					for _, key := range keys {
						Expect(cache.Contains(key)).To(BeFalse(), "memcache")
						Expect(store.Contains(key)).To(BeFalse(), "datastore")
					}
				})

				It("should return an error deleting an empty", func() {
					Expect(DeleteMulti(ctx, keys)).To(HaveOccurred())
				})
			})

			Context("when getting", func() {
				It("should get all the keys", func() {
					_, err := PutMulti(ctx, keys, &data)
					Expect(err).ToNot(HaveOccurred())

					var results = make([]string, len(keys))
					err = GetMulti(ctx, keys, &results)
					Expect(err).ToNot(HaveOccurred())
					Expect(results).To(HaveLen(len(data)))
					for i, result := range results {
						Expect(result).To(Equal(data[i]))
					}
				})

				It("should get all the data even if the memcache is empty", func() {
					_, err := PutMulti(ctx, keys, &data)
					Expect(err).ToNot(HaveOccurred())

					for _, key := range keys {
						cache.Delete(ctx, key)
					}

					var results = make([]string, len(keys))
					err = GetMulti(ctx, keys, &results)
					Expect(err).ToNot(HaveOccurred())
					Expect(results).To(HaveLen(len(data)))
					for i, result := range results {
						Expect(result).To(Equal(data[i]))
					}
				})

				It("should get all that even if it is split between the two store and memcache", func() {
					_, err := PutMulti(ctx, keys, &data)
					Expect(err).ToNot(HaveOccurred())

					for i, key := range keys {
						if i%2 == 0 {
							cache.Delete(ctx, key)
						} else {
							store.Delete(ctx, key)
						}
					}

					var results = make([]string, len(keys))
					err = GetMulti(ctx, keys, &results)
					Expect(err).ToNot(HaveOccurred())
					Expect(results).To(HaveLen(len(data)))
					for i, result := range results {
						Expect(result).To(Equal(data[i]))
					}
				})

				It("It should error on things that cannot be gotton", func() {
					var result []string
					Expect(GetMulti(ctx, keys, &result)).To(HaveOccurred())
				})
			})

			Context("when putting", func() {
				It("should put all data to the keys", func() {
					_, err := PutMulti(ctx, keys, data)
					Expect(err).ToNot(HaveOccurred())
					for _, key := range keys {
						Expect(cache.Contains(key)).To(BeTrue(), "memcache")
						Expect(store.Contains(key)).To(BeTrue(), "datastore")
					}
				})

				It("It should error on things that can be stored", func() {
					ts := []string{}
					_, err := PutMulti(ctx, keys, ts)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
