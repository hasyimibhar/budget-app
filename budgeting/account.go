package budgeting

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	Name string

	transactions []*Transaction
	closed       bool
}

func NewAccount(name string) *Account {
	return &Account{
		Name: name,

		transactions: []*Transaction{},
		closed:       false,
	}
}

func (a *Account) AddTransaction(date time.Time, amount decimal.Decimal, description string) *Transaction {
	t := NewTransaction(date, amount, description)
	a.transactions = append(a.transactions, t)
	return t
}

func (a *Account) Balance() decimal.Decimal {
	transactions := a.transactions
	sort.Sort(byDate(transactions))

	balance := decimal.New(0, -2)
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