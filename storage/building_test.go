package storage_test

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/eckyputrady/jsonapicrudexample/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// there are a lot of functions because each test can be run individually and sets up the complete
// environment. That is because we run all the specs randomized.
var _ = Describe("Building Test", func() {
	var sut *storage.BuildingStorage

	BeforeEach(func() {
		sut = storage.NewBuildingStorage()
	})

	Describe("Create", func() {
		It("Should create successfully", func() {
			sut.Insert(model.Building{})
			_, err := sut.GetOne("1")
			Expect(err).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("Should update successfully", func() {
			sut.Insert(model.Building{Address: "UG"})
			f := model.Building{Address: "G", ID: "1"}
			err := sut.Update(f)
			Expect(err).To(BeNil())
			data, _ := sut.GetOne("1")
			Expect(data).To(Equal(f))
		})

		It("Should return err if ID not found", func() {
			_, err := sut.GetOne("-1")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Delete", func() {
		It("Should delete successfully", func() {
			sut.Insert(model.Building{})
			delErr := sut.Delete("1")
			Expect(delErr).To(BeNil())
			_, err := sut.GetOne("1")
			Expect(err).ToNot(BeNil())
		})

		It("Should return err if ID not found", func() {
			delErr := sut.Delete("1")
			Expect(delErr).ToNot(BeNil())
		})
	})

	Describe("Get", func() {
		It("Should throw err if not found", func() {
			_, err := sut.GetOne("-1")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("GetAll", func() {
		It("Should return empty if no item", func() {
			data := sut.GetAll()
			Expect(data).To(BeEmpty())
		})

		It("Should return all items", func() {
			sut.Insert(model.Building{})
			sut.Insert(model.Building{})
			sut.Insert(model.Building{})
			data := sut.GetAll()
			Expect(data).To(Equal([]model.Building{
				model.Building{ID: "1"},
				model.Building{ID: "2"},
				model.Building{ID: "3"},
			}))

		})
	})

	Describe("PaginateFindAll", func() {
		BeforeEach(func() {
			sut.Insert(model.Building{Address: "A"})
			sut.Insert(model.Building{Address: "B"})
			sut.Insert(model.Building{Address: "C"})
			sut.Insert(model.Building{Address: "D"})
		})

		It("Should show empty if give out of range params", func() {
			len, data := sut.PaginatedFindAll(0, 0)
			Expect(len).To(Equal(4))
			Expect(data).To(BeEmpty())

			len, data = sut.PaginatedFindAll(10, 10)
			Expect(len).To(Equal(4))
			Expect(data).To(BeEmpty())

			len, data = sut.PaginatedFindAllLimitOffset(0, 0)
			Expect(len).To(Equal(4))
			Expect(data).To(BeEmpty())

			len, data = sut.PaginatedFindAllLimitOffset(10, 10)
			Expect(len).To(Equal(4))
			Expect(data).To(BeEmpty())
		})

		It("Should paginate correctly", func() {
			len, data := sut.PaginatedFindAll(2, 3)
			Expect(len).To(Equal(4))
			expected := []model.Building{model.Building{ID: "4", Address: "D"}}
			Expect(data).To(Equal(expected))
		})
	})

	Describe("Concurrency", func() {
		var asyncAddAndModify = func(wg *sync.WaitGroup) {
			defer wg.Done()

			// insert
			id := sut.Insert(model.Building{})

			// then either update or delete
			idInt, _ := strconv.ParseInt(id, 10, 64)
			if idInt > 50 {
				sut.Delete(id)
			} else {
				sut.Update(model.Building{ID: id, Address: "Updated"})
			}
		}

		It("Should be thread safe on writes", func() {
			var wg sync.WaitGroup
			wg.Add(100)
			for i := 0; i < 100; i++ {
				go asyncAddAndModify(&wg)
			}
			wg.Wait()

			actual := sut.GetAll()
			Expect(actual).To(HaveLen(50))

			expected := []model.Building{}
			for i := 1; i <= 50; i++ {
				expected = append(expected, model.Building{ID: fmt.Sprintf("%d", i), Address: "Updated"})
			}
			Expect(actual).To(Equal(expected))
		})
	})
})
