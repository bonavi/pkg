package errors

import "pkg/log/model"

type ErrorType struct {
	HTTPCode  int
	LogAs     LogOption
	HumanText string
}

// Error - Кастомная структура ошибки
type Error struct {

	// Тип ошибки, в который зашиты данные от разработчика
	// В случае, если этот тип снова кладется в errors.Type.Wrap, эта переменная затирается
	// Оставить первоначатльный тип ошибки можно через errors.DontEraseErrorType()
	ErrorType ErrorType `json:"-"`

	// Первоначальная ошибка. Если необходимо завернуть эту ошибку через fmt.Errorf("%w", err), то
	// Необходимо воспользоваться errors.WithAdditionalError(ErrNotFound)
	Err error `json:"-"`

	// Поскольку стандартный энкодер json в го не умеет нормально сериализовать тип ошибок, эта переменная
	// Используется для подставления значения Err прямо перед сериализацией ошибки в функции JSON
	DeveloperText string `json:"error"`

	// Человекочитаемый текст, который можно показать клиенту
	// Переменная настраивается через errors.WithCustomHumanText(messageWithFmt, args...)
	// Если значения нет, то автоматически проставляется шаблонными данными в функции middleware.DefaultErrorEncoder
	HumanText string `json:"humanText"`

	// Стектрейс от места враппинга ошибки. Если необходимо начать стектрейс с уровня выше, то
	// Необходимо воспользоваться errors.SkipThisCall(errors.<const>)
	// const = SkipThisCall - начать стектрейс на один уровень выше враппера errors.Type.Wrap по дереву
	// const = SkipPreviousCaller и остальные работают по аналогии, пропуская все больше уровней стека вызовов
	StackTrace []string `json:"path"`

	// Дополнительные параметры, направленные на дополнение ошибки контекстом, которые проставляются
	// Через errors.WithParams(key1, value1, key2, value2, ...)
	Params map[string]any `json:"parameters"`

	// Служебное поле, которое автоматически заполняется в функции middleware.DefaultErrorEncoder
	// вспомогательными данными из контекста
	//
	// Программист это поле руками не меняет!
	SystemInfo model.SystemInfo `json:"systemInfo"`
}
