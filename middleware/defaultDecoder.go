package middleware

import (
	"encoding/json"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/schema"

	"pkg/decimal"
	"pkg/errors"
)

type DecodeMethod int

const (
	DecodeSchema DecodeMethod = iota + 1
	DecodeJSON
)

var decimalConverter = func(val string) reflect.Value {
	dec, err := decimal.NewFromString(val)
	if err != nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(dec)
}

func DefaultDecoder(
	ctx *fiber.Ctx,
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
		queries := ctx.Queries()
		classicQueries := make(map[string][]string)
		for key, value := range queries {
			classicQueries[key] = []string{value}
		}
		decoder := schema.NewDecoder()
		decoder.IgnoreUnknownKeys(true)
		decoder.RegisterConverter(decimal.Decimal{}, decimalConverter)
		err = decoder.Decode(dest, classicQueries)
	case DecodeJSON:
		err = json.Unmarshal(ctx.Body(), &dest)
	}
	if err != nil {
		return errors.Default.Wrap(err).
			SkipThisCall()
	}

	return nil
}
