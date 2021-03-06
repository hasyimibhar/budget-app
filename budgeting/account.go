package budgeting

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

var (
	zero = decimal.New(0, -2)

	// ErrCannotAssignCategoryToTransfer is returned when there is an attempt
	// to create a transfer transaction without a category.
	ErrCannotAssignCategoryToTransfer = fmt.Errorf("a transfer cannot have category")
)

// Account represents a physical account which stores money
// (e.g. savings account, your wallet).
type Account struct {
	Name string

	budget              *Budget
	transactions        []*Transaction
	transactionCategory map[string][]*Transaction
	closed              bool
}

func newAccount(budget *Budget, name string, balance decimal.Decimal, date time.Time, tbb *Category) (*Account, error) {
	a := &Account{
		Name: name,

		budget:              budget,
		transactions:        []*Transaction{},
		transactionCategory: map[string][]*Transaction{},
		closed:              false,
	}

	if _, err := a.AddTransaction(date, balance, "Starting balance", tbb, nil); err != nil {
		return nil, err
	}

	return a, nil
}

// AddTransaction creates a transaction on the account.
// The rel argument is used to indicate a transfer between 2 accounts.
// If rel is not nil, a matching transaction will be created on that account
// (i.e. double entry bookkeeping).
func (a *Account) AddTransaction(
	date time.Time,
	amount decimal.Decimal,
	description string,
	category *Category,
	rel *Account) (*Transaction, error) {

	if rel != nil && category != nil {
		return nil, ErrCannotAssignCategoryToTransfer
	}

	t := newTransaction(a.budget, a, date, amount, description, category, rel)
	a.transactions = append(a.transactions, t)

	if rel != nil {
		t2 := newTransaction(a.budget, rel, date, amount.Neg(), description, category, a)
		rel.transactions = append(rel.transactions, t2)
	}

	if a.budget.earliestMonth.EarlierTime(date) {
		a.budget.earliestMonth = YearMonthFromTime(date)
	}
	if a.budget.latestMonth.LaterTime(date) {
		a.budget.latestMonth = YearMonthFromTime(date)
	}

	a.budget.setTransactionCategory(t, category)

	return t, nil
}

// Balance returns the account balance.
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
