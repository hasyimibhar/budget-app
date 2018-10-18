package budgeting

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBudget_NewBudget(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	jan := YearMonth{2018, time.January}

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	assert.True(budget.Budgeted(jan, food).Equal(dec("0.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("0.00")))
}

func TestBudget_SetBudgeted(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	jan := YearMonth{2018, time.January}
	feb := YearMonth{2018, time.February}

	budget.AddAccount("Savings", dec("100.00"), date(2018, 1, 1))
	assert.True(budget.TBB(jan).Equal(dec("100.00")))
	assert.True(budget.TBB(feb).Equal(dec("100.00"))) // TBB should carry over

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	budget.SetBudgeted(jan, food, dec("50.00"))
	budget.SetBudgeted(jan, bills, dec("12.34"))

	assert.True(budget.TBB(jan).Equal(dec("37.66")))
	assert.True(budget.TBB(feb).Equal(dec("37.66"))) // TBB should carry over

	assert.True(budget.Budgeted(jan, food).Equal(dec("50.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("12.34")))
	assert.True(budget.Budgeted(feb, food).Equal(dec("0.00")))
	assert.True(budget.Budgeted(feb, bills).Equal(dec("0.00")))
}

func TestBudget_MoveBudgeted(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	jan := YearMonth{2018, time.January}
	feb := YearMonth{2018, time.February}

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	assert.True(budget.TBB(jan).Equal(dec("0.00")))
	assert.True(budget.TBB(feb).Equal(dec("0.00")))

	budget.SetBudgeted(jan, food, dec("50.00"))
	budget.SetBudgeted(jan, bills, dec("12.34"))

	assert.True(budget.TBB(jan).Equal(dec("-62.34")))
	assert.True(budget.TBB(feb).Equal(dec("-62.34")))

	budget.MoveBudgeted(jan, food, bills, dec("10.00"))
	assert.True(budget.Budgeted(jan, food).Equal(dec("40.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("22.34")))

	budget.MoveBudgeted(jan, bills, food, dec("10.00"))
	assert.True(budget.Budgeted(jan, food).Equal(dec("50.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("12.34")))

	budget.MoveBudgeted(jan, food, bills, dec("-5.00"))
	assert.True(budget.Budgeted(jan, food).Equal(dec("55.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("7.34")))

	budget.MoveBudgeted(jan, bills, food, dec("-50.00"))
	assert.True(budget.Budgeted(jan, food).Equal(dec("5.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("57.34")))
}

func TestBudget_MoveTBB(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	jan := YearMonth{2018, time.January}

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	budget.AddAccount("Savings", dec("100.00"), date(2018, 1, 1))

	budget.MoveBudgeted(jan, budget.TBBCategory(), food, dec("15.00"))
	budget.MoveBudgeted(jan, budget.TBBCategory(), bills, dec("29.00"))
	assert.True(budget.Budgeted(jan, food).Equal(dec("15.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("29.00")))
	assert.True(budget.TBB(jan).Equal(dec("56.00")))
}

func TestBudget_TBB(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	aug := YearMonth{2018, time.August}
	sep := YearMonth{2018, time.September}
	oct := YearMonth{2018, time.October}
	nov := YearMonth{2018, time.November}

	food := budget.AddCategory("Food & Beverages")

	acc := budget.AddAccount("Savings", dec("10.00"), date(2018, 8, 1))
	acc.AddTransaction(date(2018, 10, 1), dec("100.00"), "test", budget.TBBCategory(), nil)

	assert.Equal(dec("10.00").StringFixed(2), budget.TBB(aug).StringFixed(2))
	assert.Equal(dec("10.00").StringFixed(2), budget.TBB(sep).StringFixed(2))
	assert.Equal(dec("110.00").StringFixed(2), budget.TBB(oct).StringFixed(2))
	assert.Equal(dec("110.00").StringFixed(2), budget.TBB(nov).StringFixed(2))

	budget.MoveBudgeted(sep, budget.TBBCategory(), food, dec("20.00"))

	assert.Equal(dec("0.00").StringFixed(2), budget.TBB(aug).StringFixed(2))
	assert.Equal(dec("-10.00").StringFixed(2), budget.TBB(sep).StringFixed(2))
	assert.Equal(dec("90.00").StringFixed(2), budget.TBB(oct).StringFixed(2))
	assert.Equal(dec("90.00").StringFixed(2), budget.TBB(nov).StringFixed(2))

	budget.MoveBudgeted(aug, budget.TBBCategory(), food, dec("30.00"))

	assert.Equal(dec("-20.00").StringFixed(2), budget.TBB(aug).StringFixed(2))
	assert.Equal(dec("-40.00").StringFixed(2), budget.TBB(sep).StringFixed(2))
	assert.Equal(dec("60.00").StringFixed(2), budget.TBB(oct).StringFixed(2))
	assert.Equal(dec("60.00").StringFixed(2), budget.TBB(nov).StringFixed(2))

	budget.SetBudgeted(aug, food, dec("5.00"))

	assert.Equal(dec("0.00").StringFixed(2), budget.TBB(aug).StringFixed(2))
	assert.Equal(dec("-15.00").StringFixed(2), budget.TBB(sep).StringFixed(2))
	assert.Equal(dec("85.00").StringFixed(2), budget.TBB(oct).StringFixed(2))
	assert.Equal(dec("85.00").StringFixed(2), budget.TBB(nov).StringFixed(2))
}

func TestBudget_MoveBudgetedEmpty(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	jan := YearMonth{2018, time.January}

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	budget.MoveBudgeted(jan, food, bills, dec("10.00"))
	assert.True(budget.Budgeted(jan, food).Equal(dec("-10.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("10.00")))
}
