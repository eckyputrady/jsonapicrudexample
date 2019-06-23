package resource

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/eckyputrady/jsonapicrudexample/storage"
	"github.com/manyminds/api2go"
)

// BuildingResource for api2go routes
type BuildingResource struct {
	FloorStorage    *storage.FloorStorage
	BuildingStorage *storage.BuildingStorage
}

// FindAll to satisfy api2go data source interface
func (s BuildingResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	buildings := s.BuildingStorage.GetAll()
	s.includeFloors(toRefSlice(buildings))
	return &Response{Res: buildings}, nil
}

func toRefSlice(in []model.Building) []*model.Building {
	ret := []*model.Building{}
	for _, b := range in {
		ret = append(ret, &b)
	}
	return ret
}

func (s BuildingResource) includeFloors(buildings []*model.Building) {
	for _, b := range buildings {
		b.Floors = s.FloorStorage.GetMany(b.FloorsIDs)
	}
}

func parseUintOrDefault(r api2go.Request, key string, def int) (res int, exists bool) {
	q, ok := r.QueryParams[key]
	if !ok {
		return def, false
	}

	parsed, err := strconv.ParseInt(q[0], 10, 64)
	if err != nil {
		return def, true
	}

	return int(parsed), true
}

// PaginatedFindAll can be used to load buildings in chunks
func (s BuildingResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	pageNum, pageNumExists := parseUintOrDefault(r, "page[number]", 1)
	pageSize, pageSizeExists := parseUintOrDefault(r, "page[size]", 10)
	if pageNumExists && pageSizeExists {
		n, data := s.BuildingStorage.PaginatedFindAll(pageNum, pageSize)
		s.includeFloors(toRefSlice(data))
		return uint(n), &Response{Res: data}, nil
	}

	limit, limitExists := parseUintOrDefault(r, "page[limit]", 10)
	offset, offsetExists := parseUintOrDefault(r, "page[offset]", 0)
	if limitExists && offsetExists {
		n, data := s.BuildingStorage.PaginatedFindAllLimitOffset(limit, offset)
		s.includeFloors(toRefSlice(data))
		return uint(n), &Response{Res: data}, nil
	}

	buildings := s.BuildingStorage.GetAll()
	s.includeFloors(toRefSlice(buildings))
	return uint(len(buildings)), &Response{Res: buildings}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
func (s BuildingResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	building, err := s.BuildingStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	building.Floors = s.FloorStorage.GetMany(building.FloorsIDs)

	return &Response{Res: building}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BuildingResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	building, ok := obj.(model.Building)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BuildingStorage.Insert(building)
	building.ID = id

	return &Response{Res: building, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BuildingResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BuildingStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the building
func (s BuildingResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	building, ok := obj.(model.Building)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BuildingStorage.Update(building)
	return &Response{Res: building, Code: http.StatusNoContent}, err
}
