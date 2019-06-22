package storage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/manyminds/api2go"
)

// NewBuildingStorage initializes the storage
func NewBuildingStorage() *BuildingStorage {
	return &BuildingStorage{make(map[string]*model.Building), 1}
}

// BuildingStorage stores all users
type BuildingStorage struct {
	users   map[string]*model.Building
	idCount int
}

// GetAll returns the user map (because we need the ID as key too)
func (s BuildingStorage) GetAll() map[string]*model.Building {
	return s.users
}

// GetOne user
func (s BuildingStorage) GetOne(id string) (model.Building, error) {
	user, ok := s.users[id]
	if ok {
		return *user, nil
	}
	errMessage := fmt.Sprintf("Building for id %s not found", id)
	return model.Building{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a user
func (s *BuildingStorage) Insert(c model.Building) string {
	id := fmt.Sprintf("%d", s.idCount)
	c.ID = id
	s.users[id] = &c
	s.idCount++
	return id
}

// Delete one :(
func (s *BuildingStorage) Delete(id string) error {
	_, exists := s.users[id]
	if !exists {
		return fmt.Errorf("Building with id %s does not exist", id)
	}
	delete(s.users, id)

	return nil
}

// Update a user
func (s *BuildingStorage) Update(c model.Building) error {
	_, exists := s.users[c.ID]
	if !exists {
		return fmt.Errorf("Building with id %s does not exist", c.ID)
	}
	s.users[c.ID] = &c

	return nil
}
