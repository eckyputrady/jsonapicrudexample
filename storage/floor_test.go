package storage_test

import (
	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/eckyputrady/jsonapicrudexample/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// there are a lot of functions because each test can be run individually and sets up the complete
// environment. That is because we run all the specs randomized.
var _ = Describe("Floor Test", func() {
	var sut *storage.FloorStorage

	BeforeEach(func() {
		sut = storage.NewFloorStorage()
	})

	Describe("Create", func() {
		It("Should create successfully", func() {
			sut.Insert(model.Floor{})
			_, err := sut.GetOne("1")
			Expect(err).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("Should update successfully", func() {
			sut.Insert(model.Floor{Name: "UG"})
			f := model.Floor{Name: "G", ID: "1"}
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
			sut.Insert(model.Floor{})
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
			sut.Insert(model.Floor{})
			sut.Insert(model.Floor{})
			sut.Insert(model.Floor{})
			data := sut.GetAll()
			Expect(data).To(Equal([]model.Floor{
				model.Floor{ID: "1"},
				model.Floor{ID: "2"},
				model.Floor{ID: "3"},
			}))

		})
	})

	Describe("GetMany", func() {
		It("Should get many successfully", func() {
			sut.Insert(model.Floor{})
			sut.Insert(model.Floor{})
			sut.Insert(model.Floor{})
			data := sut.GetMany([]string{"2", "3", "4"})
			Expect(data).To(HaveLen(2))
		})
	})
})
