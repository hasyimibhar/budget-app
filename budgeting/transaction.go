package budgeting

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Date        time.Time
	Description string
	Amount      decimal.Decimal

	uuid     string
	budget   *Budget
	category *Category
	rel      *Account
}

func newTransaction(
	budget *Budget,
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
		category: category,
		rel:      rel,
	}
}

func (t *Transaction) Category() *Category {
	return t.category
}
