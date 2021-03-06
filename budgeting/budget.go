package budgeting

import (
	"time"

	"github.com/shopspring/decimal"
)

type Budget struct {
	Name string

	earliestMonth YearMonth
	latestMonth   YearMonth
	tbb           *Category
	categories    map[string]*Category
	accounts      []*Account
	budgeted      map[YearMonth]monthBudget
}

type monthBudget struct {
	Month    YearMonth
	Budgeted map[string]decimal.Decimal
}

// NewBudget creates a fresh budget.
func NewBudget(name string) *Budget {
	b := &Budget{
		Name: name,

		earliestMonth: YearMonth{999999, time.December},
		latestMonth:   YearMonth{0, time.January},
		categories:    map[string]*Category{},
		accounts:      []*Account{},
		budgeted:      map[YearMonth]monthBudget{},
	}

	tbb := b.AddCategory("To Be Budgeted")
	b.tbb = tbb

	return b
}

// AddAccount creates an account within the budget.
func (b *Budget) AddAccount(name string, balance decimal.Decimal, date time.Time) *Account {
	account, _ := newAccount(b, name, balance, date, b.tbb)
	b.accounts = append(b.accounts, account)
	return account
}

// TBBCategory returns the "To Be Budgeted" category.
func (b *Budget) TBBCategory() *Category {
	return b.tbb.clone()
}

// TBB returns the "To Be Budgeted" balance for the specified month.
func (b *Budget) TBB(month YearMonth) decimal.Decimal {
	tbb := zero
	var m YearMonth

	m = month
	for {
		if b.earliestMonth.Earlier(m) {
			break
		}

		transactions := b.monthCategoryTransactions(m, b.tbb)
		for _, t := range transactions {
			tbb = tbb.Add(t.Amount)
		}

		m = m.LastMonth()
	}

	m = month
	for {
		if b.earliestMonth.Earlier(m) {
			break
		}

		budgeted := monthBudget{
			Month:    m,
			Budgeted: map[string]decimal.Decimal{},
		}

		if _, ok := b.budgeted[m]; ok {
			budgeted = b.budgeted[m]
		}

		for _, v := range budgeted.Budgeted {
			tbb = tbb.Sub(v)
		}

		m = m.LastMonth()
	}

	if tbb.GreaterThan(zero) {
		m = month.NextMonth()
		for {
			if b.latestMonth.Later(m) {
				break
			}

			budgeted := monthBudget{
				Month:    m,
				Budgeted: map[string]decimal.Decimal{},
			}

			if _, ok := b.budgeted[m]; ok {
				budgeted = b.budgeted[m]
			}

			for _, v := range budgeted.Budgeted {
				tbb = tbb.Sub(v)
				if tbb.LessThan(zero) {
					tbb = zero
					break
				}
			}

			m = m.NextMonth()
		}
	}

	return tbb
}

// AddCategory creates a budgeting category.
func (b *Budget) AddCategory(name string) *Category {
	category := newCategory(name, b)
	b.categories[category.uuid] = category
	return category
}

// Activities returns how much money has been spent for the category on the specified month.
func (b *Budget) Activities(month YearMonth, category *Category) decimal.Decimal {
	activities := zero

	transactions := b.monthCategoryTransactions(month, category)
	for _, t := range transactions {
		activities = activities.Add(t.Amount)
	}

	return activities
}

// Available returns the available budget balance for the category on the specified month.
func (b *Budget) Available(month YearMonth, category *Category) decimal.Decimal {
	available := zero

	if !b.earliestMonth.Earlier(month) {
		available = available.Add(b.Available(month.LastMonth(), category))
	}

	available = available.Add(b.Budgeted(month, category).Add(b.Activities(month, category)))

	return available
}

// Budgeted returns the budgeted amount for the category on the specified month.
func (b *Budget) Budgeted(month YearMonth, category *Category) decimal.Decimal {
	if _, ok := b.budgeted[month]; !ok {
		return zero
	}

	if category.Equal(b.tbb) {
		return b.TBB(month)
	}

	return b.budgeted[month].Budgeted[category.uuid]
}

// SetBudgeted sets the budgeted amount for the category on the specified month.
func (b *Budget) SetBudgeted(month YearMonth, category *Category, amount decimal.Decimal) {
	if category.Equal(b.tbb) {
		return
	}

	if _, ok := b.budgeted[month]; !ok {
		b.budgeted[month] = monthBudget{
			Month:    month,
			Budgeted: map[string]decimal.Decimal{},
		}
	}

	b.budgeted[month].Budgeted[category.uuid] = amount

	if b.earliestMonth.Earlier(month) {
		b.earliestMonth = month
	}
	if b.latestMonth.Later(month) {
		b.latestMonth = month
	}
}

// MoveBudgeted moves the budget balance from one category to another on the specified month.
func (b *Budget) MoveBudgeted(month YearMonth, from *Category, to *Category, amount decimal.Decimal) {
	if _, ok := b.budgeted[month]; !ok {
		b.budgeted[month] = monthBudget{
			Month:    month,
			Budgeted: map[string]decimal.Decimal{},
		}
	}

	fromAmount := b.Budgeted(month, from)
	toAmount := b.Budgeted(month, to)

	fromAmount = fromAmount.Sub(amount)
	toAmount = toAmount.Add(amount)

	b.SetBudgeted(month, from, fromAmount)
	b.SetBudgeted(month, to, toAmount)
}

func (b *Budget) setTransactionCategory(t *Transaction, c *Category) error {
	if t.Type() == TransactionTypeTransfer {
		return ErrCannotAssignCategoryToTransfer
	}

	transactions := []*Transaction{}

	if t.Category() != nil {
		transactions = t.account.transactionCategory[t.Category().uuid]
		for i, tt := range transactions {
			if tt.uuid == t.uuid {
				transactions = append(transactions[:i], transactions[i+1:]...)
				break
			}
		}

		t.account.transactionCategory[t.Category().uuid] = transactions
	}

	if c != nil {
		if _, ok := t.account.transactionCategory[c.uuid]; !ok {
			t.account.transactionCategory[c.uuid] = []*Transaction{}
		}

		t.account.transactionCategory[c.uuid] = append(t.account.transactionCategory[c.uuid], t)
	}

	t.category = c
	return nil
}

func (b *Budget) monthCategoryTransactions(month YearMonth, c *Category) []*Transaction {
	allTransactions := []*Transaction{}

	// TODO: Optimize this by using lookup table
	for _, a := range b.accounts {
		transactions := []*Transaction{}
		if _, ok := a.transactionCategory[c.uuid]; ok {
			transactions = a.transactionCategory[c.uuid]
		}

		for _, t := range transactions {
			if YearMonthFromTime(t.Date).Equal(month) {
				allTransactions = append(allTransactions, t)
			}
		}
	}

	return allTransactions
}

// YearMonth is a helper struct for representing a month of a year
// (e.g. May 2018).
type YearMonth struct {
	Year  int
	Month time.Month
}

// YearMonthFromTime creates a YearMonth from a time.Time.
func YearMonthFromTime(t time.Time) YearMonth {
	return YearMonth{
		Year:  t.Year(),
		Month: t.Month(),
	}
}

// LastMonth returns the previous month.
func (m YearMonth) LastMonth() YearMonth {
	if m.Month == time.January {
		return YearMonth{
			Year:  m.Year - 1,
			Month: time.December,
		}
	}

	return YearMonth{
		Year:  m.Year,
		Month: time.Month(int(m.Month) - 1),
	}
}

// NextMonth returns the next month.
func (m YearMonth) NextMonth() YearMonth {
	if m.Month == time.December {
		return YearMonth{
			Year:  m.Year + 1,
			Month: time.January,
		}
	}

	return YearMonth{
		Year:  m.Year,
		Month: time.Month(int(m.Month) + 1),
	}
}

// Earlier returns true if m is earlier than other.
func (m YearMonth) Earlier(other YearMonth) bool {
	if other.Year < m.Year {
		return true
	}
	if other.Year > m.Year {
		return false
	}
	return int(other.Month) < int(m.Month)
}

// EarlierTime returns true if m is earlier than t.
func (m YearMonth) EarlierTime(t time.Time) bool {
	if t.Year() < m.Year {
		return true
	}
	if t.Year() > m.Year {
		return false
	}
	return int(t.Month()) < int(m.Month)
}

// Later returns true if m is later than other.
func (m YearMonth) Later(other YearMonth) bool {
	if other.Year > m.Year {
		return true
	}
	if other.Year < m.Year {
		return false
	}
	return int(other.Month) > int(m.Month)
}

// LaterTime returns true if m is later than t.
func (m YearMonth) LaterTime(t time.Time) bool {
	if t.Year() > m.Year {
		return true
	}
	if t.Year() < m.Year {
		return false
	}
	return int(t.Month()) > int(m.Month)
}

// Equal returns true if the both years and months are equal.
func (m YearMonth) Equal(other YearMonth) bool {
	return m.Year == other.Year && m.Month == other.Month
}
