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

// FindAll floors
func (c FloorResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	buildingsID, ok := r.QueryParams["buildingsID"]
	if ok {
		buildingID := buildingsID[0]
		building, err := c.BuildingStorage.GetOne(buildingID)
		if err != nil {
			return &Response{}, err
		}

		return &Response{Res: c.FloorStorage.GetMany(building.FloorsIDs)}, nil
	}
	return &Response{Res: c.FloorStorage.GetAll()}, nil
}

// FindOne floor
func (c FloorResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.FloorStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new floor
func (c FloorResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	floor, ok := obj.(model.Floor)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.FloorStorage.Insert(floor)
	floor.ID = id
	return &Response{Res: floor, Code: http.StatusCreated}, nil
}

// Delete a floor
func (c FloorResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.FloorStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a floor
func (c FloorResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	floor, ok := obj.(model.Floor)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.FloorStorage.Update(floor)
	return &Response{Res: floor, Code: http.StatusNoContent}, err
}
