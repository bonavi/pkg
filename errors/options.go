package errors

import (
	"context"
	"fmt"
	"maps"
)

type contextKey int

const errorsParamsKey contextKey = 1

type ErrorsParams struct {
	data map[string]string
}

func GetErrorsParams(ctx context.Context) map[string]string {

	errorsParams := make(map[string]string)

	if ctx == nil {
		return errorsParams
	}

	if paramsAny := ctx.Value(errorsParamsKey); paramsAny != nil {
		if ep, ok := paramsAny.(ErrorsParams); ok {
			errorsParams = maps.Clone(ep.data)
		}
	}
	return errorsParams
}

func AddErrorsParams(ctx context.Context, paramsKV ...string) context.Context {
	if len(paramsKV) == 0 {
		return ctx
	}

	paramsCopy := maps.Clone(GetErrorsParams(ctx))

	for i := 0; i < len(paramsKV); i += 2 {
		if i+1 < len(paramsKV) {
			paramsCopy[paramsKV[i]] = paramsKV[i+1]
		} else {
			paramsCopy[paramsKV[i]] = ""
		}
	}

	return context.WithValue(ctx, errorsParamsKey, ErrorsParams{data: paramsCopy})
}

func (e Error) WithContextParams(ctx context.Context) Error {

	// Получаем параметры из контекста
	contextParams := GetErrorsParams(ctx)

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

func (e Error) WithStackTraceJump(p int) Error {

	// Если стектрейс записан и его длина больше скипа
	if e.StackTrace != nil && len(e.StackTrace) > p {

		// Удаляем первые p элементов стектрейса
		e.StackTrace = e.StackTrace[p:]
	}

	return e
}

func (e Error) WithLogOption(p LogOption) Error {

	// Меняем способ логгирования этой ошибки
	e.LogAs = p

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
