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

	tbb        *Category
	categories map[string]*Category
	accounts   []*Account
	budgeted   map[YearMonth]monthBudget
}

type monthBudget struct {
	Month    YearMonth
	Budgeted map[string]decimal.Decimal
}

func NewBudget(name string) *Budget {
	b := &Budget{
		Name: name,

		categories: map[string]*Category{},
		accounts:   []*Account{},
		budgeted:   map[YearMonth]monthBudget{},
	}

	tbb := b.AddCategory("To Be Budgeted")
	b.tbb = tbb

	return b
}

func (b *Budget) AddAccount(name string, balance decimal.Decimal, date time.Time) *Account {
	account, _ := newAccount(name, balance, date, b.tbb)
	b.accounts = append(b.accounts, account)
	return account
}

// TBB is a special category which all budget must have.
func (b *Budget) TBB() *Category {
	// TBB cannot be modified, so it's cloned to make it read-only
	return b.tbb.clone()
}

func (b *Budget) AddCategory(name string) *Category {
	category := newCategory(name)
	b.categories[category.uuid] = category
	return category
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
