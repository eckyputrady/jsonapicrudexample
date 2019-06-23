package main

import (
	"fmt"
	"net/http"

	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/eckyputrady/jsonapicrudexample/resolver"
	"github.com/eckyputrady/jsonapicrudexample/resource"
	"github.com/eckyputrady/jsonapicrudexample/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
)

func main() {
	port := 31415
	host := "localhost"
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port, Host: host})

	buildingStorage := storage.NewBuildingStorage()
	floorStorage := storage.NewFloorStorage()
	api.AddResource(model.Building{}, resource.BuildingResource{FloorStorage: floorStorage, BuildingStorage: buildingStorage})
	api.AddResource(model.Floor{}, resource.FloorResource{FloorStorage: floorStorage, BuildingStorage: buildingStorage})

	handler := api.Handler().(*httprouter.Router)
	fmt.Printf("Listening on %s:%d", host, port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
