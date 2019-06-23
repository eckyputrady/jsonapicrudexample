package resource

import (
	"errors"
	"net/http"

	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/eckyputrady/jsonapicrudexample/storage"
	"github.com/manyminds/api2go"
)

// FloorResource for api2go routes
type FloorResource struct {
	FloorStorage    *storage.FloorStorage
	BuildingStorage *storage.BuildingStorage
}

// FindAll chocolates
func (c FloorResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	usersID, ok := r.QueryParams["usersID"]
	sweets := c.FloorStorage.GetAll()
	if ok {
		// this means that we want to show all sweets of a user, this is the route
		// /v0/users/1/sweets
		userID := usersID[0]
		// filter out sweets with userID, in real world, you would just run a different database query
		filteredSweets := []model.Floor{}
		user, err := c.BuildingStorage.GetOne(userID)
		if err != nil {
			return &Response{}, err
		}
		for _, sweetID := range user.FloorsIDs {
			sweet, err := c.FloorStorage.GetOne(sweetID)
			if err != nil {
				return &Response{}, err
			}
			filteredSweets = append(filteredSweets, sweet)
		}

		return &Response{Res: filteredSweets}, nil
	}
	return &Response{Res: sweets}, nil
}

// FindOne choc
func (c FloorResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.FloorStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c FloorResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(model.Floor)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.FloorStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c FloorResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.FloorStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a choc
func (c FloorResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(model.Floor)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.FloorStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
