package uuid

import "github.com/google/uuid"

type generator struct {
	mockedValues []string
}

var defaultGenerator = new(generator)

func New() string {

	// Если нет замоканных значений, то генерируем новый UUID
	if len(defaultGenerator.mockedValues) == 0 {
		return uuid.New().String()
	}

	// Получаем первое замоканное значение
	value := defaultGenerator.mockedValues[0]

	// Удаляем его из списка замоканных значений
	defaultGenerator.mockedValues = defaultGenerator.mockedValues[1:]

	// Возвращаем его
	return value
}

func AddMockValues(mockedValues ...string) {
	defaultGenerator.mockedValues = append(defaultGenerator.mockedValues, mockedValues...)
}
