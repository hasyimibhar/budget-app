package budgeting

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type TransactionType int

const (
	// TransactionTypeIncomeExpense means the money either enters the budget (income)
	// or leaves the budget (expense).
	TransactionTypeIncomeExpense TransactionType = iota + 1

	// TransactionTypeTransfer means the money just moves between accounts within
	// the budget (e.g. drawing money from the bank account into your wallet).
	TransactionTypeTransfer
)

// Transaction represents a movement of money in the budget.
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

// Category returns the transactino category.
func (t *Transaction) Category() *Category {
	return t.category
}

// SetCategory sets the transaction category.
func (t *Transaction) SetCategory(category *Category) error {
	return t.budget.setTransactionCategory(t, category)
}

// Type returns the transaction type.
func (t *Transaction) Type() TransactionType {
	if t.rel == nil {
		return TransactionTypeIncomeExpense
	}

	return TransactionTypeTransfer
}
