package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/eckyputrady/jsonapicrudexample/model"
	"github.com/eckyputrady/jsonapicrudexample/resource"
	"github.com/eckyputrady/jsonapicrudexample/storage"
	"github.com/manyminds/api2go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// there are a lot of functions because each test can be run individually and sets up the complete
// environment. That is because we run all the specs randomized.
var _ = Describe("CrudExample", func() {
	var (
		rec *httptest.ResponseRecorder
		api *api2go.API
	)

	BeforeEach(func() {
		api = api2go.NewAPIWithBaseURL("v0", "http://localhost:31415")
		buildingStorage := storage.NewBuildingStorage()
		floorStorage := storage.NewFloorStorage()
		api.AddResource(model.Building{}, resource.BuildingResource{FloorStorage: floorStorage, BuildingStorage: buildingStorage})
		api.AddResource(model.Floor{}, resource.FloorResource{FloorStorage: floorStorage, BuildingStorage: buildingStorage})
		rec = httptest.NewRecorder()
	})

	var createBuilding = func() {
		rec = httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v0/buildings", strings.NewReader(`
		{
			"data": {
				"type": "buildings",
				"attributes": {
					"address": "Jurong East"
				}
			}
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusCreated))
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"data": {
				"id": "1",
				"type": "buildings",
				"attributes": {
					"address": "Jurong East"
				},
				"relationships": {
					"floors": {
						"data": [],
						"links": {
							"related": "http://localhost:31415/v0/buildings/1/floors",
							"self": "http://localhost:31415/v0/buildings/1/relationships/floors"
						}
					}
				}
			}
		}
		`))
	}

	It("Creates a new building", func() {
		createBuilding()
	})

	It("Gets bulding with pagination correctly", func() {
		createBuilding()

		// create another building
		rec = httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v0/buildings", strings.NewReader(`
		{
			"data": {
				"type": "buildings",
				"attributes": {
					"address": "Jurong West"
				}
			}
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusCreated))

		// Pagination
		rec = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/v0/buildings?page[limit]=1&page[offset]=1", strings.NewReader(`
		{
			"data": {
				"type": "buildings",
				"attributes": {
					"address": "Jurong East"
				}
			}
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"links": {
				"first": "http://localhost:31415/v0/buildings?page[limit]=1&page[offset]=0",
				"prev": "http://localhost:31415/v0/buildings?page[limit]=1&page[offset]=0"
			},
			"data": [
				{
					"type": "buildings",
					"id": "2",
					"attributes": {
						"address": "Jurong West"
					},
					"relationships": {
						"floors": {
							"links": {
								"related": "http://localhost:31415/v0/buildings/2/floors",
								"self": "http://localhost:31415/v0/buildings/2/relationships/floors"
							},
							"data": []
						}
					}
				}
			]
		}
		`))
	})

	var createFloor = func() {
		rec = httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v0/floors", strings.NewReader(`
		{
			"data": {
				"type": "floors",
				"attributes": {
					"name": "B2"
				}
			}
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusCreated))
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"data": {
				"id": "1",
				"type": "floors",
				"attributes": {
					"name": "B2"
				}
			}
		}
		`))
	}

	It("Creates a new floor", func() {
		createFloor()
	})

	var replaceFloors = func() {
		rec = httptest.NewRecorder()
		By("Replacing floors relationship with PATCH")

		req, err := http.NewRequest("PATCH", "/v0/buildings/1/relationships/floors", strings.NewReader(`
		{
			"data": [{
				"type": "floors",
				"id": "1"
			}]
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusNoContent))

		rec = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/v0/buildings/1", nil)
		api.Handler().ServeHTTP(rec, req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"data": {
				"attributes": {
					"address": "Jurong East"
				},
				"id": "1",
				"relationships": {
					"floors": {
						"data": [
							{
								"id": "1",
								"type": "floors"
							}
						],
						"links": {
							"related": "http://localhost:31415/v0/buildings/1/floors",
							"self": "http://localhost:31415/v0/buildings/1/relationships/floors"
						}
					}
				},
				"type": "buildings"
			},
			"included": [
				{
					"attributes": {
						"name": "B2"
					},
					"id": "1",
					"type": "floors"
				}
			]
		}
		`))
	}

	It("Creates a building with references floors", func() {
		createFloor()

		rec = httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/v0/buildings", strings.NewReader(`
    {
      "data": {
        "type": "buildings",
        "attributes": {
          "address": "Jurong East"
        },
        "relationships": {
          "floors": {
            "data": [
            {
              "id": "1",
              "type": "floors"
            }
            ]
          }
        }
      }
    }
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusCreated))
		Expect(rec.Body.String()).To(MatchJSON(`
          {
            "data": {
              "id": "1",
              "type": "buildings",
              "attributes": {
                "address": "Jurong East"
              },
              "relationships": {
                "floors": {
                  "data": [
                    {
                      "id": "1",
                      "type": "floors"
                    }
                  ],
                  "links": {
                    "related": "http://localhost:31415/v0/buildings/1/floors",
                    "self": "http://localhost:31415/v0/buildings/1/relationships/floors"
                  }
                }
              }
            }
          }
          `))
	})

	It("Replaces buildings floors", func() {
		createBuilding()
		createFloor()
		replaceFloors()
	})

	It("Deletes a buildings floor", func() {
		createBuilding()
		createFloor()
		replaceFloors()
		rec = httptest.NewRecorder()

		By("Deleting the buildings only floor with ID 1")

		req, err := http.NewRequest("DELETE", "/v0/buildings/1/relationships/floors", strings.NewReader(`
		{
			"data": [{
				"type": "floors",
				"id": "1"
			}]
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusNoContent))

		rec = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/v0/buildings/1", nil)
		api.Handler().ServeHTTP(rec, req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"data": {
				"attributes": {
					"address": "Jurong East"
				},
				"id": "1",
				"relationships": {
					"floors": {
						"data": [],
						"links": {
							"related": "http://localhost:31415/v0/buildings/1/floors",
							"self": "http://localhost:31415/v0/buildings/1/relationships/floors"
						}
					}
				},
				"type": "buildings"
			}
		}
		`))
	})

	It("Adds a buildings floor", func() {
		createBuilding()
		createFloor()
		rec = httptest.NewRecorder()

		By("Adding a floor with POST")

		req, err := http.NewRequest("POST", "/v0/buildings/1/relationships/floors", strings.NewReader(`
		{
			"data": [{
				"type": "floors",
				"id": "1"
			}]
		}
		`))
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusNoContent))

		By("Loading the building from the backend, it should have the relationship")

		rec = httptest.NewRecorder()
		req, err = http.NewRequest("GET", "/v0/buildings/1", nil)
		api.Handler().ServeHTTP(rec, req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"data": {
				"attributes": {
					"address": "Jurong East"
				},
				"id": "1",
				"relationships": {
					"floors": {
						"data": [
							{
								"id": "1",
								"type": "floors"
							}
						],
						"links": {
							"related": "http://localhost:31415/v0/buildings/1/floors",
							"self": "http://localhost:31415/v0/buildings/1/relationships/floors"
						}
					}
				},
				"type": "buildings"
			},
			"included": [
				{
					"attributes": {
						"name": "B2"
					},
					"id": "1",
					"type": "floors"
				}
			]
		}
		`))
	})

	Describe("Load floors of a building directly", func() {
		BeforeEach(func() {
			createBuilding()
			createFloor()
			replaceFloors()
			rec = httptest.NewRecorder()

			// add another floor so we have 2, only 1 is connected with the building
			req, err := http.NewRequest("POST", "/v0/floors", strings.NewReader(`
			{
				"data": {
					"type": "floors",
					"attributes": {
						"name": "G"
					}
				}
			}
			`))
			Expect(err).ToNot(HaveOccurred())
			api.Handler().ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusCreated))
			Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": {
					"id": "2",
					"type": "floors",
					"attributes": {
						"name": "G"
					}
				}
			}
			`))

			rec = httptest.NewRecorder()
		})

		It("There are 2 floors in the datastorage now", func() {
			req, err := http.NewRequest("GET", "/v0/floors", nil)
			Expect(err).ToNot(HaveOccurred())
			api.Handler().ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": [
					{
						"attributes": {
							"name": "B2"
						},
						"id": "1",
						"type": "floors"
					},
					{
						"attributes": {
							"name": "G"
						},
						"id": "2",
						"type": "floors"
					}
				]
			}
			`))
		})

		It("The building only has the previously connected floor", func() {
			req, err := http.NewRequest("GET", "/v0/buildings/1", nil)
			Expect(err).ToNot(HaveOccurred())
			api.Handler().ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": {
					"attributes": {
						"address": "Jurong East"
					},
					"id": "1",
					"relationships": {
						"floors": {
							"data": [
								{
									"id": "1",
									"type": "floors"
								}
							],
							"links": {
								"related": "http://localhost:31415/v0/buildings/1/floors",
								"self": "http://localhost:31415/v0/buildings/1/relationships/floors"
							}
						}
					},
					"type": "buildings"
				},
				"included": [
					{
						"attributes": {
							"name": "B2"
						},
						"id": "1",
						"type": "floors"
					}
				]
			}
			`))
		})

		It("Directly loading the floors", func() {
			req, err := http.NewRequest("GET", "/v0/buildings/1/floors", nil)
			Expect(err).ToNot(HaveOccurred())
			api.Handler().ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": [
				{
					"type": "floors",
					"id": "1",
					"attributes": {
						"name": "B2"
					}
				}
				]
			}
			`))
		})

		It("The relationship route works too", func() {
			req, err := http.NewRequest("GET", "/v0/buildings/1/relationships/floors", nil)
			Expect(err).ToNot(HaveOccurred())
			api.Handler().ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(MatchJSON(`
			{
				"data": [
					{
						"id": "1",
						"type": "floors"
					}
				],
				"links": {
					"related": "http://localhost:31415/v0/buildings/1/floors",
					"self": "http://localhost:31415/v0/buildings/1/relationships/floors"
				}
			}
			`))
		})
	})
})
