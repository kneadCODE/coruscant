package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Currency represents different currency types
type Currency string

const (
	// CurrencyUSD represents US Dollar
	CurrencyUSD Currency = "USD"
	// CurrencySGD represents Singapore Dollar
	CurrencySGD Currency = "SGD"
)

// Money represents a monetary amount with currency using decimal precision
type Money struct {
	Amount   decimal.Decimal
	Currency Currency
}

// NewMoney creates a new Money value object from string amount (recommended)
func NewMoney(amount string, currency Currency) (Money, error) {
	if currency == "" {
		return Money{}, fmt.Errorf("currency cannot be empty")
	}

	d, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount format: %w", err)
	}

	return Money{
		Amount:   d,
		Currency: currency,
	}, nil
}

// NewMoneyFromFloat creates Money from float64 (use with caution for financial data)
func NewMoneyFromFloat(amount float64, currency Currency) (Money, error) {
	if currency == "" {
		return Money{}, fmt.Errorf("currency cannot be empty")
	}

	d := decimal.NewFromFloat(amount)

	return Money{
		Amount:   d,
		Currency: currency,
	}, nil
}

// NewMoneyFromDecimal creates Money from existing decimal.Decimal
func NewMoneyFromDecimal(amount decimal.Decimal, currency Currency) (Money, error) {
	if currency == "" {
		return Money{}, fmt.Errorf("currency cannot be empty")
	}

	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Zero creates a zero Money value for the given currency
func Zero(currency Currency) (Money, error) {
	if currency == "" {
		return Money{}, fmt.Errorf("currency cannot be empty")
	}

	return Money{
		Amount:   decimal.Zero,
		Currency: currency,
	}, nil
}

// Float64 returns the float64 representation of the amount
func (m Money) Float64() float64 {
	amount, _ := m.Amount.Float64()
	return amount
}

// String returns formatted money string with 2 decimal places
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.Amount.StringFixed(2), m.Currency)
}

// StringFixed returns formatted money string with specified decimal places
func (m Money) StringFixed(places int32) string {
	return fmt.Sprintf("%s %s", m.Amount.StringFixed(places), m.Currency)
}

// Decimal returns the underlying decimal.Decimal
func (m Money) Decimal() decimal.Decimal {
	return m.Amount
}

// Equals checks if two Money values are equal
func (m Money) Equals(other Money) bool {
	return m.Amount.Equal(other.Amount) && m.Currency == other.Currency
}

// IsZero checks if the money amount is zero
func (m Money) IsZero() bool {
	return m.Amount.IsZero()
}

// IsPositive checks if the money amount is positive
func (m Money) IsPositive() bool {
	return m.Amount.IsPositive()
}

// IsNegative checks if the money amount is negative
func (m Money) IsNegative() bool {
	return m.Amount.IsNegative()
}

// Add adds two money amounts (must be same currency)
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", m.Currency, other.Currency)
	}

	return Money{
		Amount:   m.Amount.Add(other.Amount),
		Currency: m.Currency,
	}, nil
}

// Subtract subtracts two money amounts (must be same currency)
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("cannot subtract different currencies: %s and %s", m.Currency, other.Currency)
	}

	return Money{
		Amount:   m.Amount.Sub(other.Amount),
		Currency: m.Currency,
	}, nil
}

// Multiply multiplies the money amount by a decimal multiplier
func (m Money) Multiply(multiplier decimal.Decimal) Money {
	return Money{
		Amount:   m.Amount.Mul(multiplier),
		Currency: m.Currency,
	}
}

// MultiplyFloat multiplies the money amount by a float64 multiplier
func (m Money) MultiplyFloat(multiplier float64) Money {
	return Money{
		Amount:   m.Amount.Mul(decimal.NewFromFloat(multiplier)),
		Currency: m.Currency,
	}
}

// Divide divides the money amount by a decimal divisor
func (m Money) Divide(divisor decimal.Decimal) (Money, error) {
	if divisor.IsZero() {
		return Money{}, fmt.Errorf("cannot divide by zero")
	}

	return Money{
		Amount:   m.Amount.Div(divisor),
		Currency: m.Currency,
	}, nil
}

// Negate returns the negative of the money amount
func (m Money) Negate() Money {
	return Money{
		Amount:   m.Amount.Neg(),
		Currency: m.Currency,
	}
}

// Abs returns the absolute value of the money amount
func (m Money) Abs() Money {
	return Money{
		Amount:   m.Amount.Abs(),
		Currency: m.Currency,
	}
}

// Round rounds the money amount to the specified decimal places
func (m Money) Round(places int32) Money {
	return Money{
		Amount:   m.Amount.Round(places),
		Currency: m.Currency,
	}
}

// RoundCurrency rounds the money amount to standard currency precision
// (2 decimal places for most currencies, 0 for JPY, etc.)
func (m Money) RoundCurrency() Money {
	places := int32(2) // Default to 2 decimal places

	return m.Round(places)
}

// Compare compares two money amounts (must be same currency)
// Returns -1 if m < other, 0 if m == other, 1 if m > other
func (m Money) Compare(other Money) (int, error) {
	if m.Currency != other.Currency {
		return 0, fmt.Errorf("cannot compare different currencies: %s and %s", m.Currency, other.Currency)
	}

	return m.Amount.Cmp(other.Amount), nil
}

// GreaterThan checks if this money amount is greater than another
func (m Money) GreaterThan(other Money) (bool, error) {
	cmp, err := m.Compare(other)
	return cmp > 0, err
}

// LessThan checks if this money amount is less than another
func (m Money) LessThan(other Money) (bool, error) {
	cmp, err := m.Compare(other)
	return cmp < 0, err
}

// GreaterThanOrEqual checks if this money amount is greater than or equal to another
func (m Money) GreaterThanOrEqual(other Money) (bool, error) {
	cmp, err := m.Compare(other)
	return cmp >= 0, err
}

// LessThanOrEqual checks if this money amount is less than or equal to another
func (m Money) LessThanOrEqual(other Money) (bool, error) {
	cmp, err := m.Compare(other)
	return cmp <= 0, err
}
