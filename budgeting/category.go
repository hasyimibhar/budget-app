package budgeting

import (
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Category struct {
	Name string

	uuid   string
	budget *Budget
}

func newCategory(name string, budget *Budget) *Category {
	return &Category{
		Name:   name,
		uuid:   uuid.NewV4().String(),
		budget: budget,
	}
}

func (c *Category) Budgeted(month YearMonth) decimal.Decimal {
	return c.budget.Budgeted(month, c)
}

func (c *Category) Activities(month YearMonth) decimal.Decimal {
	return c.budget.Activities(month, c)
}

func (c *Category) Available(month YearMonth) decimal.Decimal {
	return c.budget.Available(month, c)
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
