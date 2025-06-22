package log

import (
	"context"
	"fmt"
	"maps"
	"pkg/contextMap"
	"pkg/errors"
)

func (l Log) WithContextParams(ctx context.Context) Log {

	// Получаем параметры из контекста
	contextParams := contextMap.GetMap(ctx)

	if len(contextParams) == 0 {
		return l
	}

	mergedParams := maps.Clone(l.params)
	maps.Copy(mergedParams, contextParams)
	l.params = mergedParams

	return l
}

func (l Log) WithParams(parameters ...any) Log {

	// Перебираем параметры и кладем их в мапу
	for i := 0; i < len(parameters); i += 2 {
		l.params[fmt.Sprintf("%v", parameters[i])] = fmt.Sprintf("%v", parameters[i+1])
	}

	if len(parameters)%2 != 0 {
		l.params[fmt.Sprintf("%v", parameters[len(parameters)-1])] = ""
	}

	return l
}

func (l Log) SkipThisCall() Log {

	// Если стектрейс записан и его длина больше скипа
	if l.stackTrace != nil && len(l.stackTrace) > errors.SkipThisCall {

		// Удаляем первые p элементов стектрейса
		l.stackTrace = l.stackTrace[errors.SkipThisCall:]
	}

	return l
}

func (l Log) SkipPreviousCaller() Log {

	// Если стектрейс записан и его длина больше скипа
	if l.stackTrace != nil && len(l.stackTrace) > errors.SkipPreviousCaller {

		// Удаляем первые p элементов стектрейса
		l.stackTrace = l.stackTrace[errors.SkipPreviousCaller:]
	}

	return l
}

func (l Log) Skip2PreviousCallers() Log {

	// Если стектрейс записан и его длина больше скипа
	if l.stackTrace != nil && len(l.stackTrace) > errors.Skip2PreviousCallers {

		// Удаляем первые p элементов стектрейса
		l.stackTrace = l.stackTrace[errors.Skip2PreviousCallers:]
	}

	return l
}

func WithContextParams(ctx context.Context) Log { return emptyLog().WithContextParams(ctx) }
func WithParams(parameters ...any) Log          { return emptyLog().WithParams(parameters...) }
func SkipThisCall() Log                         { return emptyLog().SkipThisCall() }
func SkipPreviousCaller() Log                   { return emptyLog().SkipPreviousCaller() }
func Skip2PreviousCallers() Log                 { return emptyLog().Skip2PreviousCallers() }
