package budgeting

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

var (
	zero = decimal.New(0, -2)

	ErrMustHaveCategory               = fmt.Errorf("an income or expense must have category")
	ErrCannotAssignCategoryToTransfer = fmt.Errorf("a transfer cannot have category")
)

type Account struct {
	Name string

	transactions []*Transaction
	closed       bool
}

func NewAccount(name string, balance decimal.Decimal, date time.Time, tbb *Category) (*Account, error) {
	a := &Account{
		Name: name,

		transactions: []*Transaction{},
		closed:       false,
	}

	if _, err := a.AddTransaction(date, balance, "Starting balance", tbb, nil); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Account) AddTransaction(
	date time.Time,
	amount decimal.Decimal,
	description string,
	category *Category,
	rel *Account) (*Transaction, error) {

	if rel == nil && category == nil {
		return nil, ErrMustHaveCategory
	} else if rel != nil && category != nil {
		return nil, ErrCannotAssignCategoryToTransfer
	}

	t := newTransaction(date, amount, description, category, rel)
	a.transactions = append(a.transactions, t)

	if rel != nil {
		t2 := newTransaction(date, amount.Neg(), description, category, a)
		rel.transactions = append(rel.transactions, t2)
	}

	return t, nil
}

func (a *Account) Balance() decimal.Decimal {
	transactions := a.transactions
	sort.Sort(byDate(transactions))

	balance := zero
	for _, t := range transactions {
		balance = balance.Add(t.Amount)
	}

	return balance
}

type byDate []*Transaction

func (p byDate) Len() int {
	return len(p)
}

func (p byDate) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p byDate) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
