package budgeting

import (
	"time"

	"github.com/shopspring/decimal"
)

type Budget struct {
	Name string

	earliestMonth   YearMonth
	latestMonth     YearMonth
	tbb             *Category
	categories      map[string]*Category
	accounts        []*Account
	budgeted        map[YearMonth]monthBudget
	tbbTransactions map[YearMonth][]*Transaction
}

type monthBudget struct {
	Month    YearMonth
	Budgeted map[string]decimal.Decimal
}

func NewBudget(name string) *Budget {
	b := &Budget{
		Name: name,

		earliestMonth:   YearMonth{999999, time.December},
		latestMonth:     YearMonth{0, time.January},
		categories:      map[string]*Category{},
		accounts:        []*Account{},
		budgeted:        map[YearMonth]monthBudget{},
		tbbTransactions: map[YearMonth][]*Transaction{},
	}

	tbb := b.AddCategory("To Be Budgeted")
	b.tbb = tbb

	return b
}

func (b *Budget) AddAccount(name string, balance decimal.Decimal, date time.Time) *Account {
	account, _ := newAccount(b, name, balance, date, b.tbb)
	b.accounts = append(b.accounts, account)
	return account
}

func (b *Budget) TBBCategory() *Category {
	return b.tbb.clone()
}

func (b *Budget) TBB(month YearMonth) decimal.Decimal {
	tbb := zero
	var m YearMonth

	m = month
	for {
		if b.earliestMonth.Earlier(m) {
			break
		}

		transactions := []*Transaction{}
		if _, ok := b.tbbTransactions[m]; ok {
			transactions = b.tbbTransactions[m]
		}

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

func (b *Budget) AddCategory(name string) *Category {
	category := newCategory(name, b)
	b.categories[category.uuid] = category
	return category
}

func (b *Budget) Activities(month YearMonth, category *Category) decimal.Decimal {
	activities := zero

	// TODO: Optimize this by using lookup table
	for _, a := range b.accounts {
		transactions := []*Transaction{}
		if _, ok := a.transactionCategory[category.uuid]; ok {
			transactions = a.transactionCategory[category.uuid]
		}

		for _, t := range transactions {
			if YearMonthFromTime(t.Date).Equal(month) {
				activities = activities.Add(t.Amount)
			}
		}
	}

	return activities
}

func (b *Budget) Available(month YearMonth, category *Category) decimal.Decimal {
	available := zero

	if !b.earliestMonth.Earlier(month) {
		available = available.Add(b.Available(month.LastMonth(), category))
	}

	available = available.Add(b.Budgeted(month, category).Add(b.Activities(month, category)))

	return available
}

func (b *Budget) Budgeted(month YearMonth, category *Category) decimal.Decimal {
	if _, ok := b.budgeted[month]; !ok {
		return zero
	}

	if category.Equal(b.tbb) {
		return b.TBB(month)
	}

	return b.budgeted[month].Budgeted[category.uuid]
}

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

func (b *Budget) addTBBTransaction(month YearMonth, t *Transaction) {
	if _, ok := b.tbbTransactions[month]; !ok {
		b.tbbTransactions[month] = []*Transaction{}
	}

	b.tbbTransactions[month] = append(b.tbbTransactions[month], t)
}

type YearMonth struct {
	Year  int
	Month time.Month
}

func YearMonthFromTime(t time.Time) YearMonth {
	return YearMonth{
		Year:  t.Year(),
		Month: t.Month(),
	}
}

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

func (m YearMonth) Earlier(other YearMonth) bool {
	if other.Year < m.Year {
		return true
	}
	if other.Year > m.Year {
		return false
	}
	return int(other.Month) < int(m.Month)
}

func (m YearMonth) EarlierTime(t time.Time) bool {
	if t.Year() < m.Year {
		return true
	}
	if t.Year() > m.Year {
		return false
	}
	return int(t.Month()) < int(m.Month)
}

func (m YearMonth) Later(other YearMonth) bool {
	if other.Year > m.Year {
		return true
	}
	if other.Year < m.Year {
		return false
	}
	return int(other.Month) > int(m.Month)
}

func (m YearMonth) LaterTime(t time.Time) bool {
	if t.Year() > m.Year {
		return true
	}
	if t.Year() < m.Year {
		return false
	}
	return int(t.Month()) > int(m.Month)
}

func (m YearMonth) Equal(other YearMonth) bool {
	return m.Year == other.Year && m.Month == other.Month
}
