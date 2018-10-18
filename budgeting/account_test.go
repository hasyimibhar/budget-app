package budgeting

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccount_NewAccount(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	assert.True(account.Balance().Equal(dec("0.00")))
}

func TestAccount_NewAccountWithBalance(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account", dec("5.00"), date(2018, 1, 1))
	assert.True(account.Balance().Equal(dec("5.00")))

	account = NewAccount("Savings Account", dec("-12.34"), date(2018, 1, 1))
	assert.True(account.Balance().Equal(dec("-12.34")))
}

func TestAccount_AddTransaction(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	tbb := NewCategory("To Be Budgeted")
	food := NewCategory("Food & Beverages")
	bills := NewCategory("Bills")

	var tr *Transaction

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("12.34"), "got some money yo", tbb, nil)
	assert.True(tr.Category().Equal(tbb))
	assert.True(account.Balance().Equal(dec("12.34")))

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("-4.11"), "hungry", food, nil)
	assert.True(tr.Category().Equal(food))
	assert.True(account.Balance().Equal(dec("8.23")))

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("10.00"), "found some money on the ground", tbb, nil)
	assert.True(tr.Category().Equal(tbb))
	assert.True(account.Balance().Equal(dec("18.23")))

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("-21.5"), "im broke", bills, nil)
	assert.True(tr.Category().Equal(bills))
	assert.True(account.Balance().Equal(dec("-3.27")))
}

func TestAccount_AddTransactionDoubleEntry(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	wallet := NewAccount("Wallet", dec("0.00"), date(2018, 1, 1))
	tbb := NewCategory("To Be Budgeted")
	food := NewCategory("Food & Beverages")

	var tr *Transaction

	account.AddTransaction(date(2018, 1, 1), dec("10.00"), "got some money", tbb, nil)
	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("-5.00"), "withdraw to wallet", nil, wallet)
	assert.Nil(tr.Category())

	assert.True(account.Balance().Equal(dec("5.00")))
	assert.True(wallet.Balance().Equal(dec("5.00")))

	wallet.AddTransaction(date(2018, 1, 1), dec("-0.5"), "buy candy", food, nil)
	tr, _ = wallet.AddTransaction(date(2018, 1, 1), dec("-3.00"), "deposit to savings", nil, account)
	assert.Nil(tr.Category())

	assert.True(account.Balance().Equal(dec("8.00")))
	assert.True(wallet.Balance().Equal(dec("1.50")))

	tr, _ = wallet.AddTransaction(date(2018, 1, 1), dec("8.00"), "withdraw again", nil, account)
	assert.Nil(tr.Category())

	assert.True(account.Balance().Equal(dec("0.00")))
	assert.True(wallet.Balance().Equal(dec("9.50")))
}

func TestAccount_AddTransactionIncomeExpenseWithoutCategory(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account", dec("0.00"), date(2018, 1, 1))

	tr, err := account.AddTransaction(date(2018, 1, 1), dec("10.00"), "got some money", nil, nil)
	assert.Nil(tr)
	assert.NotNil(err)
	assert.EqualError(err, ErrMustHaveCategory.Error())

	tr, err = account.AddTransaction(date(2018, 1, 1), dec("-3.00"), "spend some money", nil, nil)
	assert.Nil(tr)
	assert.NotNil(err)
	assert.EqualError(err, ErrMustHaveCategory.Error())
}

func TestAccount_AddTransactionTransferWithCategory(t *testing.T) {
	assert := assert.New(t)

	account := NewAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	wallet := NewAccount("Wallet", dec("0.00"), date(2018, 1, 1))
	tbb := NewCategory("To Be Budgeted")
	food := NewCategory("Food & Beverages")

	account.AddTransaction(date(2018, 1, 1), dec("10.00"), "got some money", tbb, nil)

	tr, err := account.AddTransaction(date(2018, 1, 1), dec("-5.00"), "withdraw to wallet", food, wallet)
	assert.Nil(tr)
	assert.NotNil(err)
	assert.EqualError(err, ErrCannotAssignCategoryToTransfer.Error())
}

func date(y int, m int, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func dec(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}
