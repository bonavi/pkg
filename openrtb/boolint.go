package openrtb

import (
	"encoding/json"
	"fmt"
)

// BoolInt может быть булевым или целым числом, а хранится как int
type BoolInt int

// UnmarshalJSON кастомный для обработки bool и int
func (b *BoolInt) UnmarshalJSON(data []byte) error {

	// Пробуем парсить как булевое
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err == nil {
		if boolVal {
			*b = BoolInt(1)
		} else {
			*b = BoolInt(0)
		}
		return nil
	}

	// Пробуем парсить как число
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*b = BoolInt(intVal)
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into BoolInt", string(data))
}

// MarshalJSON кастомный marshaler - всегда сериализуем как int
func (b *BoolInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*b))
}

// Int возвращает значение как int
func (b *BoolInt) Int() int {
	return int(*b)
}

// Bool возвращает значение как bool (1 = true, 0 = false)
func (b *BoolInt) Bool() bool {
	return int(*b) != 0
}
