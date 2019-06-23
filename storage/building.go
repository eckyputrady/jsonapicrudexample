package storage

import (
	"fmt"
	"sync"

	"github.com/eckyputrady/jsonapicrudexample/model"
)

// NewBuildingStorage initializes the storage
func NewBuildingStorage() *BuildingStorage {
	return &BuildingStorage{data: make(map[string]*model.Building), nextID: 1}
}

// BuildingStorage stores all buildings. This is thread-safe.
type BuildingStorage struct {
	data   map[string]*model.Building
	nextID int
	mutex  sync.RWMutex
}

// GetAll returns all buildings
func (s *BuildingStorage) GetAll() []model.Building {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := []model.Building{}
	for i := 1; i < s.nextID; i++ {
		f, err := s.GetOne(fmt.Sprintf("%d", i))
		if err != nil {
			continue
		}
		result = append(result, f)
	}

	return result
}

// PaginatedFindAll returns all buildings with pagination params
func (s *BuildingStorage) PaginatedFindAll(page int, size int) (int, []model.Building) {
	offset := size * (page - 1)
	limit := size
	return s.PaginatedFindAllLimitOffset(limit, offset)
}

// PaginatedFindAllLimitOffset returns all building with paginated params limit & offset
func (s *BuildingStorage) PaginatedFindAllLimitOffset(limit int, offset int) (int, []model.Building) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// normalize
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}

	result := []model.Building{}
	all := s.GetAll()
	len := len(all)
	for i := offset; i < offset+limit && i < len; i++ {
		result = append(result, all[i])
	}

	return len, result
}

// GetOne user
func (s *BuildingStorage) GetOne(id string) (model.Building, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, exists := s.data[id]
	if !exists {
		return model.Building{}, fmt.Errorf("Building with id %s does not exist", id)
	}

	return *data, nil
}

// Insert a user
func (s *BuildingStorage) Insert(c model.Building) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	c.ID = fmt.Sprintf("%d", s.nextID)
	s.data[c.ID] = &c
	s.nextID++
	return c.ID
}

// Delete one building
func (s *BuildingStorage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.data[id]
	if !exists {
		return fmt.Errorf("Building with id %s does not exist", id)
	}
	delete(s.data, id)

	return nil
}

// Update a building
func (s *BuildingStorage) Update(c model.Building) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.data[c.ID]
	if !exists {
		return fmt.Errorf("Building with id %s does not exist", c.ID)
	}
	s.data[c.ID] = &c

	return nil
}
