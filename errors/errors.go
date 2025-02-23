package errors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"

	"pkg/stackTrace"
)

// New создает новую ошибку
func (et ErrorType) New(ctx context.Context, msg string, opts ...Option) error {

	// Получаем опции в готовой структуре
	options := mergeOptions(opts...)

	// Получаем глубину пути согласно переданным опциям
	skip := stackTrace.ThisCall
	if options.stackTrace != nil {
		skip = *options.stackTrace
	}

	// Создаем новую ошибку
	customErr := Error{
		SystemInfo:    settings.SystemInfo,
		ErrorType:     et,
		Ctx:           ctx,
		DeveloperText: "",
		HumanText:     options.HumanText,
		Err:           errors.New(msg),
		StackTrace:    stackTrace.GetStackTrace(skip + 1),
		Params:        options.params,
		UserInfo:      ctx.Value(settings.UserInfoContextKey),
		LogAs:         et.logOptionByDefault(),
	}

	// Если передан тип логирования, то добавляем его
	if options.logAs != nil {
		customErr.LogAs = *options.logAs
	}

	return customErr
}

// Wrap оборачивает ошибку
func (et ErrorType) Wrap(ctx context.Context, err error, opts ...Option) error {

	// Получаем опции в готовой структуре
	options := mergeOptions(opts...)

	// Определяем глубину пути по умолчанию
	skip := stackTrace.ThisCall

	var customErr Error

	// Если это уже обернутая ошибка
	if As(err, &customErr) {

		// Если передан текст для пользователя, то затираем его
		if options.HumanText != "" {
			customErr.HumanText = options.HumanText
		}

		// Если передана глубина пути, то используем ее
		if options.stackTrace != nil {
			skip = *options.stackTrace
		}
		customErr.StackTrace = stackTrace.GetStackTrace(skip + 1)

		// Если передан текст ошибки, то добавляем его к исходной ошибке
		if options.errorf != nil {
			customErr.Err = fmt.Errorf("%w: %w", *options.errorf, customErr.Err)
		}

		// Если передан параметр, что тип ошибки затирать не надо, то не затираем
		if options.dontEraseErrorType == nil {
			customErr.ErrorType = et
		}

	} else { // Если это не обернутая ошибка

		// Если передана глубина пути, то используем ее
		if options.stackTrace != nil {
			skip = *options.stackTrace
		}

		// Если это не обернутая ошибка, то создаем новую
		customErr = Error{
			SystemInfo:    settings.SystemInfo,
			ErrorType:     et,
			Ctx:           ctx,
			DeveloperText: "",
			HumanText:     options.HumanText,
			Err:           err,
			StackTrace:    stackTrace.GetStackTrace(skip + 1),
			Params:        options.params,
			UserInfo:      ctx.Value(settings.UserInfoContextKey),
			LogAs:         et.logOptionByDefault(),
		}

		if options.errorf != nil {
			customErr.Err = fmt.Errorf("%w: %w", *options.errorf, err)
		}
	}

	// Добавляем параметры
	if customErr.Params == nil {
		customErr.Params = make(map[string]string)
		maps.Copy(customErr.Params, options.params)
	}

	if options.logAs != nil {
		customErr.LogAs = *options.logAs
	}

	return customErr
}

// JSON преобразует ошибку в json
func JSON(err Error) ([]byte, error) {

	// Подставляем значение ошибки в текстовую переменную DeveloperTextError, поскольку сериализатор не умеет
	// нормально обрабатывать тип error
	err.DeveloperText = err.Err.Error()
	byt, e := json.Marshal(err)
	if e != nil {
		return nil, InternalServer.Wrap(err.Ctx, e)
	}

	return byt, nil
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

func New(err string) error {
	return errors.New(err)
}
