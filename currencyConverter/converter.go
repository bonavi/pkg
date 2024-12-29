package currencyConverter

import (
	"pkg/decimal"

	"pkg/errors"
)

const DefaultCurrency = "RUB"

func Convert(price decimal.Decimal, fromCurrency, toCurrency string, rates map[string]decimal.Decimal) (decimal.Decimal, error) {
	// Проверяем, существуют ли курсы для обеих валют в мапе.
	fromRate, ok1 := rates[fromCurrency]
	toRate, ok2 := rates[toCurrency]
	if !ok1 || !ok2 {
		return decimal.Decimal{}, errors.InternalServer.New("Exchange rate for one of the currencies not found",
			errors.ParamsOption(
				"fromCurrency", fromCurrency,
				"toCurrency", toCurrency,
			),
			errors.SkipThisCallOption(),
		)
	}

	// Конвертируем цену в базовую валюту (USD)
	basePrice := price.Mul(fromRate)

	// Конвертируем базовую валюту в целевую валюту.
	convertedPrice := basePrice.Div(toRate)

	return convertedPrice, nil
}

func Coefficient(fromCurrency, toCurrency string, rates map[string]decimal.Decimal) (decimal.Decimal, error) {

	// Проверяем, существуют ли курсы для обеих валют в мапе.
	fromRate, ok1 := rates[fromCurrency]
	toRate, ok2 := rates[toCurrency]
	if !ok1 || !ok2 {
		return decimal.Decimal{}, errors.InternalServer.New("Exchange rate for one of the currencies not found",
			errors.ParamsOption(
				"fromCurrency", fromCurrency,
				"toCurrency", toCurrency,
			),
			errors.SkipThisCallOption(),
		)
	}

	// Конвертируем курс валюты в коэффициент
	coefficientFromTo := fromRate.Div(toRate)

	return coefficientFromTo, nil
}

func ConvertWithCoefficient(price decimal.Decimal, coefficient decimal.Decimal) decimal.Decimal {

	// Конвертируем цену в целевую валюту с учетом коэффициента
	convertedPrice := price.Mul(coefficient)

	return convertedPrice
}
