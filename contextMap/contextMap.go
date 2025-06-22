package contextMap

import (
	"context"
	"fmt"
	"maps"
)

type contextKey int

const ErrorsParamsKey contextKey = 1

type ErrorsParams struct {
	data map[string]any
}

func NewContextMap(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := ctx.Value(ErrorsParamsKey).(ErrorsParams); !ok {
		ctx = context.WithValue(ctx, ErrorsParamsKey, ErrorsParams{data: make(map[string]any)})
	}

	return ctx
}

func GetMap(ctx context.Context) map[string]any {

	errorsParams := make(map[string]any)

	if ctx == nil {
		return errorsParams
	}

	if paramsAny := ctx.Value(ErrorsParamsKey); paramsAny != nil {
		if ep, ok := paramsAny.(ErrorsParams); ok {
			errorsParams = maps.Clone(ep.data)
		}
	}
	return errorsParams
}

func AddValue(ctx context.Context, paramsKV ...any) context.Context {
	if len(paramsKV) == 0 {
		return ctx
	}

	paramsCopy := GetMap(ctx)

	for i := 0; i < len(paramsKV); i += 2 {

		key, ok := paramsKV[i].(string)
		if !ok {
			key = fmt.Sprintf("%+v", paramsKV[i]) // Приводим к строке
		}

		if i+1 < len(paramsKV) {
			paramsCopy[key] = paramsKV[i+1]
		} else {
			paramsCopy[key] = ""
		}
	}

	return context.WithValue(ctx, ErrorsParamsKey, ErrorsParams{data: paramsCopy})
}

func GetValue(ctx context.Context, key string) (any, bool) {
	if ctx == nil {
		return nil, false
	}

	params := GetMap(ctx)
	value, exists := params[key]
	return value, exists
}

func RemoveValue(ctx context.Context, keys ...string) context.Context {
	if len(keys) == 0 {
		return ctx
	}

	paramsCopy := GetMap(ctx)
	for _, k := range keys {
		delete(paramsCopy, k)
	}

	return context.WithValue(ctx, ErrorsParamsKey, ErrorsParams{data: paramsCopy})
}

func Join(src context.Context, dest context.Context) context.Context {
	if src == nil {
		return dest
	}

	if dest == nil {
		context.Background()
	}

	srcParams := GetMap(src)
	destParams := GetMap(dest)

	maps.Copy(destParams, srcParams)

	return context.WithValue(dest, ErrorsParamsKey, ErrorsParams{data: destParams})
}
