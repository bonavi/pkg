package decoder

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/schema"

	"pkg/decimal"
	"pkg/errors"
	"pkg/reflectUtils"
)

type DecodeMethod int

const (
	DecodeSchema DecodeMethod = iota + 1
	DecodeJSON
)

func Decode(
	ctx context.Context,
	r *http.Request,
	dest any,
	decodeSchemas ...DecodeMethod,
) (err error) {

	// Проверяем типы данных
	if err = reflectUtils.CheckPointerToStruct(ctx, dest); err != nil {
		return err
	}

	// Проходимся по каждому
	for _, decodeSchema := range decodeSchemas {
		switch decodeSchema {
		case DecodeSchema:
			err = schema.NewDecoder().Decode(dest, r.URL.Query())
		case DecodeJSON:
			err = json.NewDecoder(r.Body).Decode(dest)
		default:
			break
		}
		if err != nil {
			return errors.BadRequest.Wrap(ctx, err,
				errors.SkipThisCallOption(),
			)
		}
	}

	return nil
}

var decimalConverter = func(val string) reflect.Value {
	dec, err := decimal.NewFromString(val)
	if err != nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(dec)
}

func DecodeFiber(
	ctx context.Context,
	c *fiber.Ctx,
	decodeSchema DecodeMethod,
	dest any,
) (err error) {

	// TODO тут не учтено, что если в JSON приходит массив объектов,
	// 	тогда в dest будет слайс структур, что для этой проверки будет являться указателем и вызывает ошибку
	// Проверяем типы данных
	// if err = reflectUtils.CheckPointerToStruct(dest); err != nil {
	//	 return err
	// }

	switch decodeSchema {
	case DecodeSchema:
		queries := c.Queries()
		classicQueries := make(map[string][]string)
		for key, value := range queries {
			classicQueries[key] = []string{value}
		}
		decoder := schema.NewDecoder()
		decoder.RegisterConverter(decimal.Decimal{}, decimalConverter) //nolint:exhaustruct
		err = decoder.Decode(dest, classicQueries)
	case DecodeJSON:
		err = json.Unmarshal(c.Body(), &dest)
	}
	if err != nil {
		return errors.BadRequest.Wrap(ctx, err,
			errors.SkipThisCallOption(),
		)
	}

	return nil
}
