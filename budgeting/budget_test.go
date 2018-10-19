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
	assert.Equal(dec("100.00").StringFixed(2), budget.TBB(jan).StringFixed(2))
	assert.Equal(dec("100.00").StringFixed(2), budget.TBB(feb).StringFixed(2)) // TBB should carry over

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	budget.SetBudgeted(jan, food, dec("50.00"))
	budget.SetBudgeted(jan, bills, dec("12.34"))

	assert.Equal(dec("37.66").StringFixed(2), budget.TBB(jan).StringFixed(2))
	assert.Equal(dec("37.66").StringFixed(2), budget.TBB(feb).StringFixed(2)) // TBB should carry over

	assert.Equal(dec("50.00").StringFixed(2), budget.Budgeted(jan, food).StringFixed(2))
	assert.Equal(dec("12.34").StringFixed(2), budget.Budgeted(jan, bills).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), budget.Budgeted(feb, food).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), budget.Budgeted(feb, bills).StringFixed(2))
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

	budget.SetBudgeted(aug, food, dec("30.00"))
	budget.SetBudgeted(oct, food, dec("20.00"))
	acc.AddTransaction(date(2018, 9, 1), dec("5.00"), "test", budget.TBBCategory(), nil)

	assert.Equal(dec("-20.00").StringFixed(2), budget.TBB(aug).StringFixed(2))
	assert.Equal(dec("-35.00").StringFixed(2), budget.TBB(sep).StringFixed(2))
	assert.Equal(dec("45.00").StringFixed(2), budget.TBB(oct).StringFixed(2))
	assert.Equal(dec("45.00").StringFixed(2), budget.TBB(nov).StringFixed(2))
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

func TestBudget_Activities(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	dcm := YearMonth{2017, time.December}
	jan := YearMonth{2018, time.January}
	feb := YearMonth{2018, time.February}
	mar := YearMonth{2018, time.March}

	acc := budget.AddAccount("Savings", dec("100.00"), date(2018, 1, 1))
	food := budget.AddCategory("Food")
	bills := budget.AddCategory("Bills")

	budget.SetBudgeted(jan, food, dec("50.00"))
	budget.SetBudgeted(jan, bills, dec("50.00"))

	acc.AddTransaction(date(2018, 1, 2), dec("-5.00"), "lunch", food, nil)
	acc.AddTransaction(date(2018, 1, 5), dec("-3.00"), "dinner", food, nil)
	acc.AddTransaction(date(2018, 1, 2), dec("-23.00"), "electricity", bills, nil)
	acc.AddTransaction(date(2018, 1, 3), dec("-10.00"), "topup phone", bills, nil)

	acc.AddTransaction(date(2018, 2, 10), dec("-15.00"), "fancy dinner", food, nil)
	acc.AddTransaction(date(2018, 2, 2), dec("-4.00"), "water bill", bills, nil)

	t1, _ := acc.AddTransaction(date(2018, 1, 2), dec("-10.00"), "some stuff", nil, nil)
	t2, _ := acc.AddTransaction(date(2018, 2, 2), dec("-20.00"), "some other stuff", nil, nil)

	assert.Equal(dec("50.00").StringFixed(2), food.Budgeted(jan).StringFixed(2))
	assert.Equal(dec("-8.00").StringFixed(2), food.Activities(jan).StringFixed(2))
	assert.Equal(dec("42.00").StringFixed(2), food.Available(jan).StringFixed(2))

	assert.Equal(dec("50.00").StringFixed(2), bills.Budgeted(jan).StringFixed(2))
	assert.Equal(dec("-33.00").StringFixed(2), bills.Activities(jan).StringFixed(2))
	assert.Equal(dec("17.00").StringFixed(2), bills.Available(jan).StringFixed(2))

	assert.Equal(dec("0.00").StringFixed(2), food.Budgeted(feb).StringFixed(2))
	assert.Equal(dec("-15.00").StringFixed(2), food.Activities(feb).StringFixed(2))
	assert.Equal(dec("27.00").StringFixed(2), food.Available(feb).StringFixed(2))

	assert.Equal(dec("0.00").StringFixed(2), bills.Budgeted(feb).StringFixed(2))
	assert.Equal(dec("-4.00").StringFixed(2), bills.Activities(feb).StringFixed(2))
	assert.Equal(dec("13.00").StringFixed(2), bills.Available(feb).StringFixed(2))

	assert.Equal(dec("0.00").StringFixed(2), food.Budgeted(dcm).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), food.Activities(dcm).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), food.Available(dcm).StringFixed(2))

	assert.Equal(dec("0.00").StringFixed(2), bills.Budgeted(dcm).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), bills.Activities(dcm).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), bills.Available(dcm).StringFixed(2))

	assert.Equal(dec("0.00").StringFixed(2), food.Budgeted(mar).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), food.Activities(mar).StringFixed(2))
	assert.Equal(dec("27.00").StringFixed(2), food.Available(mar).StringFixed(2))

	assert.Equal(dec("0.00").StringFixed(2), bills.Budgeted(mar).StringFixed(2))
	assert.Equal(dec("0.00").StringFixed(2), bills.Activities(mar).StringFixed(2))
	assert.Equal(dec("13.00").StringFixed(2), bills.Available(mar).StringFixed(2))

	t1.SetCategory(food)
	t2.SetCategory(bills)

	assert.Equal(dec("-18.00").StringFixed(2), food.Activities(jan).StringFixed(2))
	assert.Equal(dec("32.00").StringFixed(2), food.Available(jan).StringFixed(2))

	assert.Equal(dec("-24.00").StringFixed(2), bills.Activities(feb).StringFixed(2))
	assert.Equal(dec("-7.00").StringFixed(2), bills.Available(feb).StringFixed(2))

	t1.SetCategory(nil)
	t2.SetCategory(nil)

	assert.Equal(dec("-8.00").StringFixed(2), food.Activities(jan).StringFixed(2))
	assert.Equal(dec("42.00").StringFixed(2), food.Available(jan).StringFixed(2))

	assert.Equal(dec("-4.00").StringFixed(2), bills.Activities(feb).StringFixed(2))
	assert.Equal(dec("13.00").StringFixed(2), bills.Available(feb).StringFixed(2))
}
