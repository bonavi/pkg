package errors

import (
	"context"
	"fmt"
	"maps"
	"pkg/contextMap"
)

func (e Error) WithContextParams(ctx context.Context) Error {

	// Получаем параметры из контекста
	contextParams := contextMap.GetMap(ctx)

	if len(contextParams) == 0 {
		return e
	}

	mergedParams := maps.Clone(e.Params)
	maps.Copy(mergedParams, contextParams)
	e.Params = mergedParams

	return e
}

func (e Error) WithParams(parameters ...any) Error {

	// Перебираем параметры и кладем их в мапу
	for i := 0; i < len(parameters); i += 2 {
		e.Params[fmt.Sprintf("%v", parameters[i])] = fmt.Sprintf("%v", parameters[i+1])
	}

	if len(parameters)%2 != 0 {
		e.Params[fmt.Sprintf("%v", parameters[len(parameters)-1])] = ""
	}

	return e
}

func (e Error) SkipThisCall() Error {

	// Если стектрейс записан и его длина больше скипа
	if e.StackTrace != nil && len(e.StackTrace) > SkipThisCall-1 {

		// Удаляем первые p элементов стектрейса
		e.StackTrace = e.StackTrace[SkipThisCall-1:]
	}

	return e
}

func (e Error) SkipPreviousCaller() Error {

	// Если стектрейс записан и его длина больше скипа
	if e.StackTrace != nil && len(e.StackTrace) > SkipPreviousCaller-1 {

		// Удаляем первые p элементов стектрейса
		e.StackTrace = e.StackTrace[SkipPreviousCaller-1:]
	}

	return e
}

func (e Error) Skip2PreviousCallers() Error {

	// Если стектрейс записан и его длина больше скипа
	if e.StackTrace != nil && len(e.StackTrace) > Skip2PreviousCallers-1 {

		// Удаляем первые p элементов стектрейса
		e.StackTrace = e.StackTrace[Skip2PreviousCallers-1:]
	}

	return e
}

func (e Error) WithLogOption(p LogOption) Error {

	// Меняем способ логгирования этой ошибки
	e.ErrorType.LogAs = p

	return e
}

func (e Error) WithCustomHumanText(p string, args ...any) Error {

	// Форматируем и помещаем в HumanText
	e.HumanText = fmt.Sprintf(p, args...)

	return e
}

func (e Error) WithAdditionalError(err error) Error {

	// Добавляем дополнительную ошибку к этой через дефолтную обертку
	e.Err = fmt.Errorf("%w: %w", e.Err, err)

	return e
}
