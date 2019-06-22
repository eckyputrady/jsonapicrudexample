/*
Package examples shows how to implement a basic CRUD for two data structures with the api2go server functionality.
To play with this example server you can run some of the following curl requests

In order to demonstrate dynamic baseurl handling for requests, apply the --header="REQUEST_URI:https://www.your.domain.example.com" parameter to any of the commands.

Create a new building:
	curl -X POST http://localhost:31415/v0/buildings -d '{"data" : {"type" : "buildings" , "attributes": {"building-name" : "marvin"}}}'

List buildings:
	curl -X GET http://localhost:31415/v0/buildings

List paginated buildings:
	curl -X GET 'http://localhost:31415/v0/buildings?page\[offset\]=0&page\[limit\]=2'
OR
	curl -X GET 'http://localhost:31415/v0/buildings?page\[number\]=1&page\[size\]=2'

Update:
	curl -vX PATCH http://localhost:31415/v0/buildings/1 -d '{ "data" : {"type" : "buildings", "id": "1", "attributes": {"building-name" : "better marvin"}}}'

Delete:
	curl -vX DELETE http://localhost:31415/v0/buildings/2

Create a floorolate with the name sweet
	curl -X POST http://localhost:31415/v0/floorolates -d '{"data" : {"type" : "floorolates" , "attributes": {"name" : "Ritter Sport", "taste": "Very Good"}}}'

Create a building with a sweet
	curl -X POST http://localhost:31415/v0/buildings -d '{"data" : {"type" : "buildings" , "attributes": {"building-name" : "marvin"}, "relationships": {"sweets": {"data": [{"type": "floorolates", "id": "1"}]}}}}'

List a buildings sweets
	curl -X GET http://localhost:31415/v0/buildings/1/sweets

Replace a buildings sweets
	curl -X PATCH http://localhost:31415/v0/buildings/1/relationships/sweets -d '{"data" : [{"type": "floorolates", "id": "2"}]}'

Add a sweet
	curl -X POST http://localhost:31415/v0/buildings/1/relationships/sweets -d '{"data" : [{"type": "floorolates", "id": "2"}]}'

Remove a sweet
	curl -X DELETE http://localhost:31415/v0/buildings/1/relationships/sweets -d '{"data" : [{"type": "floorolates", "id": "2"}]}'
*/
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
	api := api2go.NewAPIWithResolver("v0", &resolver.RequestURL{Port: port})
	buildingStorage := storage.NewBuildingStorage()
	floorStorage := storage.NewFloorStorage()
	api.AddResource(model.Building{}, resource.BuildingResource{FloorStorage: floorStorage, BuildingStorage: buildingStorage})
	api.AddResource(model.Floor{}, resource.FloorResource{FloorStorage: floorStorage, BuildingStorage: buildingStorage})

	fmt.Printf("Listening on :%d", port)
	handler := api.Handler().(*httprouter.Router)
	// It is also possible to get the instance of julienschmidt/httprouter and add more custom routes!
	handler.GET("/hello-world", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello World!\n")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
