# JSONAPI CRUD Example

The code is heavily inspired from https://github.com/manyminds/api2go/tree/master/examples

## APIs

Create a new building:
	curl -X POST http://localhost:31415/v0/buildings -d '{"data" : {"type" : "buildings" , "attributes": {"address" : "hello"}}}'

List buildings:
	curl -X GET http://localhost:31415/v0/buildings

List paginated buildings:
	curl -X GET 'http://localhost:31415/v0/buildings?page\[offset\]=0&page\[limit\]=2'
OR
	curl -X GET 'http://localhost:31415/v0/buildings?page\[number\]=1&page\[size\]=2'

Update:
	curl -vX PATCH http://localhost:31415/v0/buildings/1 -d '{ "data" : {"type" : "buildings", "id": "1", "attributes": {"address" : "hello 2"}}}'

Delete:
	curl -vX DELETE http://localhost:31415/v0/buildings/2

Create a floor with the name "UG"
	curl -X POST http://localhost:31415/v0/floors -d '{"data" : {"type" : "floors" , "attributes": {"name" : "UG", "taste": "Very Good"}}}'

Create a building with a floor
	curl -X POST http://localhost:31415/v0/buildings -d '{"data" : {"type" : "buildings" , "attributes": {"address" : "hello"}, "relationships": {"floors": {"data": [{"type": "floors", "id": "1"}]}}}}'

List a buildings floors
	curl -X GET http://localhost:31415/v0/buildings/1/floors

Replace a buildings floors
	curl -X PATCH http://localhost:31415/v0/buildings/1/relationships/floors -d '{"data" : [{"type": "floors", "id": "2"}]}'

Add a floor
	curl -X POST http://localhost:31415/v0/buildings/1/relationships/floors -d '{"data" : [{"type": "floors", "id": "2"}]}'

Remove a floor
	curl -X DELETE http://localhost:31415/v0/buildings/1/relationships/floors -d '{"data" : [{"type": "floors", "id": "2"}]}'