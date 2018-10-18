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

	uuid string
	rel  *Account
}

func newTransaction(date time.Time, amount decimal.Decimal, description string, rel *Account) *Transaction {
	return &Transaction{
		Date:        date,
		Amount:      amount,
		Description: description,

		uuid: uuid.NewV4().String(),
		rel:  rel,
	}
}
