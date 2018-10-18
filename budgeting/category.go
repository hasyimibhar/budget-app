package budgeting

import (
	uuid "github.com/satori/go.uuid"
)

type Category struct {
	Name string

	uuid string
}

// ToBeBudgeted is a special category that all budget
// will have.
var ToBeBudgeted = &Category{
	Name: "To Be Budgeted",
	uuid: "tbb",
}

func NewCategory(name string) *Category {
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
