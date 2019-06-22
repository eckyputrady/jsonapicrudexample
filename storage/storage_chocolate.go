package storage

import (
	"fmt"
	"sort"

	"github.com/eckyputrady/jsonapicrudexample/model"
)

// sorting
type byID []model.Floor

func (c byID) Len() int {
	return len(c)
}

func (c byID) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c byID) Less(i, j int) bool {
	return c[i].GetID() < c[j].GetID()
}

// NewFloorStorage initializes the storage
func NewFloorStorage() *FloorStorage {
	return &FloorStorage{make(map[string]*model.Floor), 1}
}

// FloorStorage stores all of the tasty chocolate, needs to be injected into
// User and Floor Resource. In the real world, you would use a database for that.
type FloorStorage struct {
	floors  map[string]*model.Floor
	idCount int
}

// GetAll of the chocolate
func (s FloorStorage) GetAll() []model.Floor {
	result := []model.Floor{}
	for key := range s.floors {
		result = append(result, *s.floors[key])
	}

	sort.Sort(byID(result))
	return result
}

// GetOne tasty chocolate
func (s FloorStorage) GetOne(id string) (model.Floor, error) {
	choc, ok := s.floors[id]
	if ok {
		return *choc, nil
	}

	return model.Floor{}, fmt.Errorf("Floor for id %s not found", id)
}

// Insert a fresh one
func (s *FloorStorage) Insert(c model.Floor) string {
	id := fmt.Sprintf("%d", s.idCount)
	c.ID = id
	s.floors[id] = &c
	s.idCount++
	return id
}

// Delete one :(
func (s *FloorStorage) Delete(id string) error {
	_, exists := s.floors[id]
	if !exists {
		return fmt.Errorf("Floor with id %s does not exist", id)
	}
	delete(s.floors, id)

	return nil
}

// Update updates an existing chocolate
func (s *FloorStorage) Update(c model.Floor) error {
	_, exists := s.floors[c.ID]
	if !exists {
		return fmt.Errorf("Floor with id %s does not exist", c.ID)
	}
	s.floors[c.ID] = &c

	return nil
}
