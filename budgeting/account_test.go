package budgeting

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccount_NewAccount(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account")
	assert.True(account.Balance().Equal(dec("0.00")))
}

func TestAccount_AddTransaction(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account")

	account.AddTransaction(date(2018, 1, 1), dec("12.34"), "test")
	assert.True(account.Balance().Equal(dec("12.34")))

	account.AddTransaction(date(2018, 1, 1), dec("-4.11"), "test")
	assert.True(account.Balance().Equal(dec("8.23")))

	account.AddTransaction(date(2018, 1, 1), dec("10.00"), "test")
	assert.True(account.Balance().Equal(dec("18.23")))

	account.AddTransaction(date(2018, 1, 1), dec("-21.5"), "test")
	assert.True(account.Balance().Equal(dec("-3.27")))
}

func date(y int, m int, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func dec(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}
