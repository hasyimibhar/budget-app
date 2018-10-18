package budgeting

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var month = YearMonth{2018, time.January}

func TestAccount_NewAccount(t *testing.T) {
	assert := assert.New(t)
	b := NewBudget("My Budget")

	account := b.AddAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	assert.True(account.Balance().Equal(dec("0.00")))
	assert.True(b.TBB(month).Equal(dec("0.00")))
}

func TestAccount_NewAccountWithBalance(t *testing.T) {
	assert := assert.New(t)
	b := NewBudget("My Budget")

	a1 := b.AddAccount("Savings Account 1", dec("5.00"), date(2018, 1, 1))
	assert.True(a1.Balance().Equal(dec("5.00")))
	assert.True(b.TBB(month).Equal(dec("5.00")))

	a2 := b.AddAccount("Savings Account 2", dec("-12.34"), date(2018, 1, 1))
	assert.True(a2.Balance().Equal(dec("-12.34")))
	assert.True(b.TBB(month).Equal(dec("-7.34")))
}

func TestAccount_AddTransaction(t *testing.T) {
	assert := assert.New(t)
	b := NewBudget("My Budget")

	account := b.AddAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	food := b.AddCategory("Food & Beverages")
	bills := b.AddCategory("Bills")

	var tr *Transaction

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("12.34"), "got some money yo", b.TBBCategory(), nil)
	assert.True(tr.Category().Equal(b.TBBCategory()))
	assert.True(account.Balance().Equal(dec("12.34")))

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("-4.11"), "hungry", food, nil)
	assert.True(tr.Category().Equal(food))
	assert.True(account.Balance().Equal(dec("8.23")))

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("10.00"), "found some money on the ground", b.TBBCategory(), nil)
	assert.True(tr.Category().Equal(b.TBBCategory()))
	assert.True(account.Balance().Equal(dec("18.23")))

	tr, _ = account.AddTransaction(date(2018, 1, 1), dec("-21.5"), "im broke", bills, nil)
	assert.True(tr.Category().Equal(bills))
	assert.True(account.Balance().Equal(dec("-3.27")))
}

func TestAccount_AddTransactionDoubleEntry(t *testing.T) {
	assert := assert.New(t)
	b := NewBudget("My Budget")

	account := b.AddAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	wallet := b.AddAccount("Wallet", dec("0.00"), date(2018, 1, 1))
	food := b.AddCategory("Food & Beverages")

	var tr *Transaction

	account.AddTransaction(date(2018, 1, 1), dec("10.00"), "got some money", b.TBBCategory(), nil)
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
	b := NewBudget("My Budget")

	account := b.AddAccount("Savings Account", dec("0.00"), date(2018, 1, 1))

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
	b := NewBudget("My Budget")

	account := b.AddAccount("Savings Account", dec("0.00"), date(2018, 1, 1))
	wallet := b.AddAccount("Wallet", dec("0.00"), date(2018, 1, 1))
	food := b.AddCategory("Food & Beverages")

	account.AddTransaction(date(2018, 1, 1), dec("10.00"), "got some money", b.TBBCategory(), nil)

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
