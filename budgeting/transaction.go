package budgeting

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type TransactionType int

const (
	TransactionTypeIncomeExpense TransactionType = iota + 1
	TransactionTypeTransfer
)

type Transaction struct {
	Date        time.Time
	Description string
	Amount      decimal.Decimal

	uuid     string
	budget   *Budget
	account  *Account
	category *Category
	rel      *Account
}

func newTransaction(
	budget *Budget,
	account *Account,
	date time.Time,
	amount decimal.Decimal,
	description string,
	category *Category,
	rel *Account) *Transaction {

	return &Transaction{
		Date:        date,
		Amount:      amount,
		Description: description,

		uuid:     uuid.NewV4().String(),
		budget:   budget,
		account:  account,
		category: category,
		rel:      rel,
	}
}

func (t *Transaction) Category() *Category {
	return t.category
}

func (t *Transaction) SetCategory(category *Category) error {
	return t.budget.setTransactionCategory(t, category)
}

func (t *Transaction) Type() TransactionType {
	if t.rel == nil {
		return TransactionTypeIncomeExpense
	}

	return TransactionTypeTransfer
}
