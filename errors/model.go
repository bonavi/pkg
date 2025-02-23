package errors

import (
	"context"
)

// Error - Кастомная структура ошибки
type Error struct {

	// Тип ошибки, в который зашит HTTP-код
	// В случае, если этот тип снова кладется в errors.Type.Wrap, эта переменная затирается
	// Оставить первоначатльный тип ошибки можно через errors.DontEraseErrorType()
	ErrorType ErrorType `json:"-"`

	// Первоначальная ошибка. Если необходимо завернуть эту ошибку через fmt.Errorf("%w", err), то
	// Необходимо воспользоваться errors.ErrorfOption("additionalInfo: %w")
	Err error `json:"-"`

	Ctx context.Context `json:"-"`

	// Поскольку стандартный энкодер json в го не умеет нормально сериализовать тип ошибок, эта переменная
	// Используется для подставления значения Err прямо перед сериализацией ошибки в функции JSON
	DeveloperText string `json:"error"`

	// Человекочитаемый текст, который можно показать клиенту
	// Переменная настраивается через errors.HumanTextOption(messageWithFmt, args...)
	// Если значения нет, то автоматически проставляется шаблонными данными в функции middleware.DefaultErrorEncoder
	HumanText string `json:"humanText"`

	// Стектрейс от места враппинга ошибки. Если необходимо начать стектрейс с уровня выше, то
	// Необходимо воспользоваться errors.SkipThisCallOption(errors.<const>)
	// const = SkipThisCall - начать стектрейс на один уровень выше враппера errors.Type.Wrap по дереву
	// const = SkipPreviousCaller и остальные работают по аналогии, пропуская все больше уровней стека вызовов
	StackTrace []string `json:"path"`

	// Дополнительные параметры, направленные на дополнение ошибки контекстом, которые проставляются
	// Через errors.ParamsOption(key1, value1, key2, value2, ...)
	Params map[string]string `json:"parameters,omitempty"`

	// Служебное поле, которое автоматически заполняется в функции middleware.DefaultErrorEncoder
	// вспомогательными данными из контекста
	UserInfo any `json:"userInfo,omitempty"`

	// Служебное поле, которое автоматически заполняется в функции middleware.DefaultErrorEncoder
	// вспомогательными данными из контекста
	SystemInfo any `json:"systemInfo,omitempty"`

	// Параметр, определяющий уровень логгирования ошибки в функции middleware.DefaultErrorEncoder
	// Настраивается через errors.LogAsOption(LogOption)
	LogAs LogOption `json:"-"`
}

// Error реализует протокол ошибок для использования нашей структуры в качестве error параметра
func (err Error) Error() string {
	return err.Err.Error()
}

// CastError приводит приедшую ошибку к нашей кастомной ошибке, если пришедшая ошибка не кастомная
// То оборачиает ее и добавляет данные о том, что ошибка не обернута
func CastError(ctx context.Context, err error) Error {
	var customErr Error
	if !As(err, &customErr) {
		err = InternalServer.Wrap(ctx, err,
			SkipThisCallOption(),
			ParamsOption("error", "Ошибка не обернута, путь неверный"),
		)
		_ = As(err, &customErr)
	}
	return customErr
}
