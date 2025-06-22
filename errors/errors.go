package errors

import (
	"context"
	"errors"
	"net/http"

	"pkg/log/model"
	"pkg/stackTrace"
)

const (
	SkipThisCall = iota + 2
	SkipPreviousCaller
	Skip2PreviousCallers
)

// Error реализует протокол ошибок для использования нашей структуры в качестве error параметра
func (e Error) Error() string {
	return e.Err.Error()
}

// LogOption - Перечисление, необходимое для конкретизации уровня логгирования ошибки
type LogOption int

const (
	LogAsError LogOption = iota
	LogAsWarning
	LogAsDebug
	LogAsInfo
	LogNone
)

// New создает новую ошибку
func (typ ErrorType) New(msg string) Error {

	var systemInfo model.SystemInfo

	// Создаем новую ошибку
	return Error{
		ErrorType:     typ,                                    // Меняется при повторном оборачивании через Wrap
		DeveloperText: "",                                     // Служебное поле, используется для сериализации в JSON
		HumanText:     "",                                     // Проставляется в хэндлере на дефолтное значение или на значение из опциональной функции WithCustomHumanText
		Err:           errors.New(msg),                        // Исходная ошибка, можно добавить дополнительную ошибку через WithAdditionalError для проведения логики через Is
		StackTrace:    stackTrace.GetStackTrace(SkipThisCall), // По дефолту получаем стектрейс от места создания этой ошибки, если необходимо урезать часть системных вызовов, можно использовать WithStackTraceJump
		Params:        make(map[string]any),                   // Дополнительные параметры, проставляются через WithParams или забираются из контекста через WithContextParams
		SystemInfo:    systemInfo,                             // Проставляется в обработчике ошибок эндпоинта
	}
}

// Wrap оборачивает ошибку
func (typ ErrorType) Wrap(err error) Error {

	var customErr Error
	var systemInfo model.SystemInfo

	if As(err, &customErr) { // Если это уже обернутая ошибка

		// Возвращаем ее
		return customErr

	} else { // Если это не обернутая ошибка

		// Если это не обернутая ошибка, то создаем новую
		return Error{
			ErrorType:     typ,                                    // Меняется при повторном оборачивании через Wrap
			DeveloperText: "",                                     // Служебное поле, используется для сериализации в JSON
			HumanText:     "",                                     // Проставляется в хэндлере на дефолтное значение или на значение из опциональной функции WithCustomHumanText
			Err:           err,                                    // Исходная ошибка, можно добавить дополнительную ошибку через WithAdditionalError для проведения логики через Is
			StackTrace:    stackTrace.GetStackTrace(SkipThisCall), // По дефолту получаем стектрейс от места создания этой ошибки, если необходимо урезать часть системных вызовов, можно использовать SkipThisCall
			Params:        make(map[string]any),                   // Дополнительные параметры, проставляются через WithParams или забираются из контекста через WithContextParams
			SystemInfo:    systemInfo,                             // Проставляется в обработчике ошибок эндпоинта
		}
	}
}

// CastError приводит приедшую ошибку к нашей кастомной ошибке, если пришедшая ошибка не кастомная
// То оборачиает ее и добавляет данные о том, что ошибка не обернута
func CastError(err error) Error {
	var customErr Error
	if !As(err, &customErr) {
		err = Default.Wrap(err).SkipThisCall().WithParams("error", "Ошибка не обернута, путь неверный")
		_ = As(err, &customErr)
	}
	return customErr
}

// As используется для вызова стандартной функции As
func As(get error, target any) bool {
	return errors.As(get, target)
}

// Unwrap используется для разворачивания завернутых с помощью fmt.Errorf("%w", err) ошибок
// default(default(1)) -> default(1)
// custom(default(default(1))) -> custom(default(1))
func Unwrap(err error) error {
	var customErr Error
	if As(err, &customErr) {
		customErr.Err = errors.Unwrap(customErr.Err)
		return customErr
	} else {
		return errors.Unwrap(err)
	}
}

// Is используется для проверки типов любой комбинации дефолтных и кастомных ошибок
func Is(err error, target error) bool {

	var customErr, customTarget Error
	if As(err, &customErr) {
		if As(target, &customTarget) {
			return errors.Is(customErr.Err, customTarget.Err) // custom - custom
		} else {
			return errors.Is(customErr.Err, target) // custom - default
		}
	} else {
		if As(target, &customTarget) {
			return errors.Is(err, customTarget.Err) // default - custom
		} else {
			return errors.Is(err, target) // default - default
		}
	}
}

func IsContextError(err error) bool {
	return Is(err, context.Canceled) || Is(err, context.DeadlineExceeded)
}

func New(text string) error {
	return errors.New(text)
}

var Default = ErrorType{
	HTTPCode:  http.StatusInternalServerError,
	LogAs:     LogAsError,
	HumanText: "",
}
