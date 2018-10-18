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

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	budget.SetBudgeted(jan, food, dec("50.00"))
	budget.SetBudgeted(jan, bills, dec("12.34"))

	assert.True(budget.Budgeted(jan, food).Equal(dec("50.00")))
	assert.True(budget.Budgeted(jan, bills).Equal(dec("12.34")))
	assert.True(budget.Budgeted(feb, food).Equal(dec("0.00")))
	assert.True(budget.Budgeted(feb, bills).Equal(dec("0.00")))
}

func TestBudget_MoveBudgeted(t *testing.T) {
	assert := assert.New(t)

	budget := NewBudget("My Budget")
	jan := YearMonth{2018, time.January}

	food := budget.AddCategory("Food & Beverages")
	bills := budget.AddCategory("Bills")

	budget.SetBudgeted(jan, food, dec("50.00"))
	budget.SetBudgeted(jan, bills, dec("12.34"))

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
