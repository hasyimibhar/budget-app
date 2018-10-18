package budgeting

import (
	uuid "github.com/satori/go.uuid"
)

type Category struct {
	Name string

	uuid string
}

func newCategory(name string) *Category {
	return &Category{
		Name: name,
		uuid: uuid.NewV4().String(),
	}
}

func (c *Category) Equal(other *Category) bool {
	if other == nil {
		return false
	}

	return c.uuid == other.uuid
}

func (c *Category) clone() *Category {
	return &Category{
		Name: c.Name,
		uuid: c.uuid,
	}
}
