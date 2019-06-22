package model

// Floor of the building
type Floor struct {
	ID   string `json:"-"`
	Name string `json:"name"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Floor) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Floor) SetID(id string) error {
	c.ID = id
	return nil
}
