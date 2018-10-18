package budgeting

import (
	"time"

	"github.com/shopspring/decimal"
)

type YearMonth struct {
	Year  int
	Month time.Month
}

type Budget struct {
	Name string

	accounts []*Account
	budgeted map[YearMonth]monthBudget
}

type monthBudget struct {
	Month    YearMonth
	Budgeted map[string]decimal.Decimal
}

func NewBudget(name string) *Budget {
	return &Budget{
		Name: name,

		accounts: []*Account{},
		budgeted: map[YearMonth]monthBudget{},
	}
}

func (b *Budget) Budgeted(month YearMonth, category *Category) decimal.Decimal {
	if _, ok := b.budgeted[month]; !ok {
		return zero
	}

	return b.budgeted[month].Budgeted[category.uuid]
}

func (b *Budget) SetBudgeted(month YearMonth, category *Category, amount decimal.Decimal) {
	if _, ok := b.budgeted[month]; !ok {
		b.budgeted[month] = monthBudget{
			Month:    month,
			Budgeted: map[string]decimal.Decimal{},
		}
	}

	b.budgeted[month].Budgeted[category.uuid] = amount
}

func (b *Budget) MoveBudgeted(month YearMonth, from *Category, to *Category, amount decimal.Decimal) {
	if _, ok := b.budgeted[month]; !ok {
		b.budgeted[month] = monthBudget{
			Month:    month,
			Budgeted: map[string]decimal.Decimal{},
		}
	}

	fromAmount := zero
	if _, ok := b.budgeted[month].Budgeted[from.uuid]; ok {
		fromAmount = b.budgeted[month].Budgeted[from.uuid]
	}

	toAmount := zero
	if _, ok := b.budgeted[month].Budgeted[to.uuid]; ok {
		toAmount = b.budgeted[month].Budgeted[to.uuid]
	}

	fromAmount = fromAmount.Sub(amount)
	toAmount = toAmount.Add(amount)

	b.budgeted[month].Budgeted[from.uuid] = fromAmount
	b.budgeted[month].Budgeted[to.uuid] = toAmount
}
