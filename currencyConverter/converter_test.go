package currencyConverter

import (
	"testing"

	"pkg/decimal"
	"pkg/testUtils"
)

func TestConvert(t *testing.T) {
	rates := map[string]decimal.Decimal{
		"USD": decimal.NewFromFloat(1),
		"RUB": decimal.NewFromFloat(0.01),
		"EUR": decimal.NewFromFloat(1.1),
	}
	type args struct {
		price        decimal.Decimal
		fromCurrency string
		toCurrency   string
	}
	tests := []struct {
		name    string
		args    args
		want    decimal.Decimal
		wantErr error
	}{
		{
			name: "1. Convert 1 USD to RUB",
			args: args{
				price:        decimal.NewFromFloat(1),
				fromCurrency: "USD",
				toCurrency:   "RUB",
			},
			want:    decimal.NewFromFloat(100),
			wantErr: nil,
		},
		{
			name: "2. Convert 1 RUB to USD",
			args: args{
				price:        decimal.NewFromFloat(1),
				fromCurrency: "RUB",
				toCurrency:   "USD",
			},
			want:    decimal.NewFromFloat(0.01),
			wantErr: nil,
		},
		{
			name: "3. Convert 1 USD to USD",
			args: args{
				price:        decimal.NewFromFloat(1),
				fromCurrency: "USD",
				toCurrency:   "USD",
			},
			want:    decimal.NewFromFloat(1),
			wantErr: nil,
		},
		{
			name: "4. Convert 0.5 EUR to RUB",
			args: args{
				price:        decimal.NewFromFloat(0.5),
				fromCurrency: "EUR",
				toCurrency:   "RUB",
			},
			want:    decimal.NewFromFloat(55),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.args.price, tt.args.fromCurrency, tt.args.toCurrency, rates)
			testUtils.CheckError(t, tt.wantErr, err, false)
			if !got.Equal(tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
