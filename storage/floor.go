package storage

import (
	"fmt"
	"sync"

	"github.com/eckyputrady/jsonapicrudexample/model"
)

// FloorStorage stores all floors. This is thread-safe.
type FloorStorage struct {
	data   map[string]*model.Floor
	nextID int
	mutex  sync.RWMutex
}

// NewFloorStorage initializes the storage
func NewFloorStorage() *FloorStorage {
	return &FloorStorage{data: make(map[string]*model.Floor), nextID: 1}
}

// GetAll of the chocolate
func (s *FloorStorage) GetAll() []model.Floor {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := []model.Floor{}
	for i := 1; i < s.nextID; i++ {
		f, err := s.GetOne(fmt.Sprintf("%d", i))
		if err != nil {
			continue
		}
		result = append(result, f)
	}

	return result
}

// GetOne floor
func (s *FloorStorage) GetOne(id string) (model.Floor, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, exists := s.data[id]
	if !exists {
		return model.Floor{}, fmt.Errorf("Floor with id %s does not exist", id)
	}

	return *data, nil
}

// GetMany floors by IDs
func (s *FloorStorage) GetMany(ids []string) []model.Floor {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := []model.Floor{}
	for _, id := range ids {
		f, err := s.GetOne(id)
		if err != nil {
			continue
		}
		result = append(result, f)
	}

	return result
}

// Insert a fresh one
func (s *FloorStorage) Insert(c model.Floor) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	c.ID = fmt.Sprintf("%d", s.nextID)
	s.data[c.ID] = &c
	s.nextID++
	return c.ID
}

// Delete one floor
func (s *FloorStorage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.data[id]
	if !exists {
		return fmt.Errorf("Floor with id %s does not exist", id)
	}
	delete(s.data, id)

	return nil
}

// Update an existing floor
func (s *FloorStorage) Update(c model.Floor) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.data[c.ID]
	if !exists {
		return fmt.Errorf("Floor with id %s does not exist", c.ID)
	}
	s.data[c.ID] = &c

	return nil
}
