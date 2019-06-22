package resource

import (
	"errors"
	"net/http"
	"sort"
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
	var result []model.Building
	users := s.BuildingStorage.GetAll()

	for _, user := range users {
		// get all sweets for the user
		user.Floors = []*model.Floor{}
		for _, chocolateID := range user.FloorsIDs {
			choc, err := s.FloorStorage.GetOne(chocolateID)
			if err != nil {
				return &Response{}, err
			}
			user.Floors = append(user.Floors, &choc)
		}
		result = append(result, *user)
	}

	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load users in chunks
func (s BuildingResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []model.Building
		number, size, offset, limit string
		keys                        []int
	)
	users := s.BuildingStorage.GetAll()

	for k := range users {
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		keys = append(keys, int(i))
	}
	sort.Ints(keys)

	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}

	if size != "" {
		sizeI, err := strconv.ParseUint(size, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		numberI, err := strconv.ParseUint(number, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		start := sizeI * (numberI - 1)
		for i := start; i < start+sizeI; i++ {
			if i >= uint64(len(users)) {
				break
			}
			result = append(result, *users[strconv.FormatInt(int64(keys[i]), 10)])
		}
	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		for i := offsetI; i < offsetI+limitI; i++ {
			if i >= uint64(len(users)) {
				break
			}
			result = append(result, *users[strconv.FormatInt(int64(keys[i]), 10)])
		}
	}

	return uint(len(users)), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BuildingResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	user, err := s.BuildingStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	user.Floors = []*model.Floor{}
	for _, chocolateID := range user.FloorsIDs {
		choc, err := s.FloorStorage.GetOne(chocolateID)
		if err != nil {
			return &Response{}, err
		}
		user.Floors = append(user.Floors, &choc)
	}
	return &Response{Res: user}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BuildingResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(model.Building)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BuildingStorage.Insert(user)
	user.ID = id

	return &Response{Res: user, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BuildingResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BuildingStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the user
func (s BuildingResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(model.Building)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BuildingStorage.Update(user)
	return &Response{Res: user, Code: http.StatusNoContent}, err
}
