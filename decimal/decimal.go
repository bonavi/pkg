package decimal

import (
	"github.com/shopspring/decimal"
)

type Decimal struct {
	Decimal decimal.Decimal
}

var Zero = Decimal{decimal.Zero}

func OffQuotesInJSON() {
	decimal.MarshalJSONWithoutQuotes = true
}

func (d Decimal) IsZero() bool {
	return d.Decimal.IsZero()
}

const normalizeCoef = 10
const precision = 16

func Normalize(d decimal.Decimal) decimal.Decimal {
	return d.
		Mul(decimal.NewFromInt(normalizeCoef)).
		Div(decimal.NewFromInt(normalizeCoef))
}

func NewFromString(s string) (Decimal, error) {
	d, err := decimal.NewFromString(s)
	return Decimal{d}, err
}

func NewFromInt(i int) Decimal {
	d := decimal.NewFromInt(int64(i))
	return Decimal{d}
}

func NewFromFloat(i float64) Decimal {
	d := decimal.NewFromFloat(i)
	return Decimal{d}
}

func (d Decimal) Mul(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Mul(d2.Decimal)}
}

func (d Decimal) Add(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Add(d2.Decimal)}
}

func (d Decimal) Sub(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Sub(d2.Decimal)}
}

func (d Decimal) Div(d2 Decimal) Decimal {
	return Decimal{d.Decimal.Div(d2.Decimal)}
}

func (d Decimal) LessThan(d2 Decimal) bool {
	return d.Decimal.LessThan(d2.Decimal)
}

func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	return d.Decimal.LessThanOrEqual(d2.Decimal)
}

func (d Decimal) Equal(d2 Decimal) bool {
	return d.Decimal.Equal(d2.Decimal)
}

func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.Decimal.GreaterThan(d2.Decimal)
}

func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	return d.Decimal.GreaterThanOrEqual(d2.Decimal)
}

func (d Decimal) String() string {
	return d.Decimal.String()
}

func (d Decimal) Float64() float64 {
	f, _ := d.Decimal.Float64()
	return f
}

func (d Decimal) DeepEqual(d2 Decimal) bool {
	return d.Decimal.Equal(d2.Decimal)
}

func (d Decimal) Round() Decimal {
	return Decimal{d.Decimal.Round(precision)}
}
