package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name        string
		amount      string
		currency    Currency
		wantErr     bool
		expectedF64 float64
	}{
		{"valid USD", "123.45", CurrencyUSD, false, 123.45},
		{"zero amount", "0.00", CurrencyUSD, false, 0.0},
		{"negative amount", "-50.25", CurrencySGD, false, -50.25},
		{"empty currency", "100.00", "", true, 0},
		{"invalid amount", "invalid", CurrencyUSD, true, 0},
		{"high precision", "123.456789", CurrencyUSD, false, 123.456789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.currency, money.Currency)
				assert.InDelta(t, tt.expectedF64, money.Float64(), 0.000001)
			}
		})
	}
}

func TestNewMoneyFromFloat(t *testing.T) {
	money, err := NewMoneyFromFloat(123.45, CurrencyUSD)
	require.NoError(t, err)

	assert.Equal(t, CurrencyUSD, money.Currency)
	assert.InDelta(t, 123.45, money.Float64(), 0.01)

	// Test empty currency
	_, err = NewMoneyFromFloat(100.0, "")
	assert.Error(t, err)
}

func TestZero(t *testing.T) {
	money, err := Zero(CurrencyUSD)
	require.NoError(t, err)

	assert.Equal(t, CurrencyUSD, money.Currency)
	assert.True(t, money.IsZero())
	assert.Equal(t, 0.0, money.Float64())

	// Test empty currency
	_, err = Zero("")
	assert.Error(t, err)
}

func TestMoney_String(t *testing.T) {
	money, _ := NewMoney("123.45", CurrencyUSD)
	assert.Equal(t, "123.45 USD", money.String())

	money, _ = NewMoney("-50.00", CurrencySGD)
	assert.Equal(t, "-50.00 SGD", money.String())

	// Test high precision
	money, _ = NewMoney("123.456789", CurrencyUSD)
	assert.Equal(t, "123.46 USD", money.String()) // Should round to 2 places for display

	// Test StringFixed
	assert.Equal(t, "123.4568 USD", money.StringFixed(4))
}

func TestMoney_Checks(t *testing.T) {
	zero, _ := NewMoney("0.00", CurrencyUSD)
	positive, _ := NewMoney("100.00", CurrencyUSD)
	negative, _ := NewMoney("-50.00", CurrencyUSD)

	assert.True(t, zero.IsZero())
	assert.False(t, positive.IsZero())
	assert.False(t, negative.IsZero())

	assert.False(t, zero.IsPositive())
	assert.True(t, positive.IsPositive())
	assert.False(t, negative.IsPositive())

	assert.False(t, zero.IsNegative())
	assert.False(t, positive.IsNegative())
	assert.True(t, negative.IsNegative())
}

func TestMoney_Add(t *testing.T) {
	money1, _ := NewMoney("100.00", CurrencyUSD)
	money2, _ := NewMoney("50.50", CurrencyUSD)

	result, err := money1.Add(money2)
	require.NoError(t, err)

	expected, _ := NewMoney("150.50", CurrencyUSD)
	assert.True(t, result.Equals(expected))

	// Test different currencies
	money3, _ := NewMoney("100.00", CurrencySGD)
	_, err = money1.Add(money3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add different currencies")
}

func TestMoney_Subtract(t *testing.T) {
	money1, _ := NewMoney("100.00", CurrencyUSD)
	money2, _ := NewMoney("30.50", CurrencyUSD)

	result, err := money1.Subtract(money2)
	require.NoError(t, err)

	expected, _ := NewMoney("69.50", CurrencyUSD)
	assert.True(t, result.Equals(expected))

	// Test different currencies
	money3, _ := NewMoney("100.00", CurrencySGD)
	_, err = money1.Subtract(money3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot subtract different currencies")
}

func TestMoney_Multiply(t *testing.T) {
	money, _ := NewMoney("100.00", CurrencyUSD)

	// Test decimal multiplication
	result := money.MultiplyFloat(1.5)
	expected, _ := NewMoney("150.00", CurrencyUSD)
	assert.True(t, result.Equals(expected))

	// Test high precision
	precise, _ := NewMoney("10.333333", CurrencyUSD)
	result2 := precise.MultiplyFloat(3.0)
	assert.InDelta(t, 30.999999, result2.Float64(), 0.000001)
}

func TestMoney_Negate(t *testing.T) {
	positive, _ := NewMoney("100.00", CurrencyUSD)
	negative, _ := NewMoney("-100.00", CurrencyUSD)

	assert.True(t, positive.Negate().Equals(negative))
	assert.True(t, negative.Negate().Equals(positive))
}

func TestMoney_Abs(t *testing.T) {
	positive, _ := NewMoney("100.00", CurrencyUSD)
	negative, _ := NewMoney("-100.00", CurrencyUSD)

	assert.True(t, positive.Abs().Equals(positive))
	assert.True(t, negative.Abs().Equals(positive))
}

func TestMoney_Round(t *testing.T) {
	money, _ := NewMoney("123.456789", CurrencyUSD)

	rounded2 := money.Round(2)
	expected2, _ := NewMoney("123.46", CurrencyUSD)
	assert.True(t, rounded2.Equals(expected2))

	rounded0 := money.Round(0)
	expected0, _ := NewMoney("123", CurrencyUSD)
	assert.True(t, rounded0.Equals(expected0))
}

func TestMoney_RoundCurrency(t *testing.T) {
	// Test USD (2 decimal places)
	usd, _ := NewMoney("123.456", CurrencyUSD)
	roundedUSD := usd.RoundCurrency()
	expectedUSD, _ := NewMoney("123.46", CurrencyUSD)
	assert.True(t, roundedUSD.Equals(expectedUSD))
}

func TestMoney_Compare(t *testing.T) {
	money1, _ := NewMoney("100.00", CurrencyUSD)
	money2, _ := NewMoney("100.00", CurrencyUSD)
	money3, _ := NewMoney("150.00", CurrencyUSD)
	money4, _ := NewMoney("50.00", CurrencyUSD)

	// Test equal
	cmp, err := money1.Compare(money2)
	assert.NoError(t, err)
	assert.Equal(t, 0, cmp)

	// Test greater than
	cmp, err = money3.Compare(money1)
	assert.NoError(t, err)
	assert.Equal(t, 1, cmp)

	greater, err := money3.GreaterThan(money1)
	assert.NoError(t, err)
	assert.True(t, greater)

	// Test less than
	cmp, err = money4.Compare(money1)
	assert.NoError(t, err)
	assert.Equal(t, -1, cmp)

	less, err := money4.LessThan(money1)
	assert.NoError(t, err)
	assert.True(t, less)

	// Test different currencies
	eurMoney, _ := NewMoney("100.00", CurrencySGD)
	_, err = money1.Compare(eurMoney)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot compare different currencies")
}

func TestMoney_Equals(t *testing.T) {
	money1, _ := NewMoney("100.00", CurrencyUSD)
	money2, _ := NewMoney("100.00", CurrencyUSD)
	money3, _ := NewMoney("100.00", CurrencySGD)
	money4, _ := NewMoney("200.00", CurrencyUSD)

	assert.True(t, money1.Equals(money2))
	assert.False(t, money1.Equals(money3)) // Different currency
	assert.False(t, money1.Equals(money4)) // Different amount
}
