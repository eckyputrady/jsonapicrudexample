package model

import (
	"errors"

	"github.com/manyminds/api2go/jsonapi"
)

// Building represents a building
type Building struct {
	ID string `json:"-"`
	//rename the username field to user-name.
	Address   string   `json:"address"`
	Floors    []Floor  `json:"-"`
	FloorsIDs []string `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Building) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Building) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Building) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "floors",
			Name: "floors",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Building) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, chocolateID := range u.FloorsIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   chocolateID,
			Type: "floors",
			Name: "floors",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Building) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Floors {
		result = append(result, u.Floors[key])
	}

	return result
}

// SetToManyReferenceIDs sets the floors reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Building) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "floors" {
		u.FloorsIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new floors
func (u *Building) AddToManyIDs(name string, IDs []string) error {
	if name == "floors" {
		u.FloorsIDs = append(u.FloorsIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some floors
func (u *Building) DeleteToManyIDs(name string, IDs []string) error {
	if name == "floors" {
		for _, ID := range IDs {
			for pos, oldID := range u.FloorsIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.FloorsIDs = append(u.FloorsIDs[:pos], u.FloorsIDs[pos+1:]...)
				}
			}
		}
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}
