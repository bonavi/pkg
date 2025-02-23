package validator

import (
	"context"
	"reflect"

	"pkg/errors"
	"pkg/stackTrace"
)

type validatorProtocol interface {
	Validate(context.Context) error
}

// Validate валидирует полученную структуру по тегам в декларативном формате
func Validate(ctx context.Context, data any) error {

	// Валидируем структуру
	if err := ZeroValue(ctx, data); err != nil {
		return errors.BadRequest.Wrap(ctx, err,
			errors.SkipThisCallOption(),
		)
	}

	// Если структура реализует интерфейс валидатора, то валидируем ее с помощью функции
	if v, ok := data.(validatorProtocol); ok {
		if err := v.Validate(ctx); err != nil {
			return errors.BadRequest.Wrap(ctx, err,
				errors.SkipThisCallOption(),
			)
		}
	}

	return nil
}

func zeroValue(ctx context.Context, requestStruct any, tag string, depth int) (tags []string, err error) {

	reflectValue := reflect.ValueOf(requestStruct)

	// Если передана структура, а не указатель на структуру, приводим к указателю на структуру
	switch {
	case reflectValue.Kind() == reflect.Struct:
		reflectValue = reflect.New(reflectValue.Type()).Elem()
		reflectValue.Set(reflect.ValueOf(requestStruct))
	case reflectValue.Kind() == reflect.Ptr && reflectValue.Elem().Kind() == reflect.Struct:
		// Если передан указатель на структуру, разыменовываем
		reflectValue = reflectValue.Elem()
	default:
		return tags, errors.InternalServer.New(ctx, "Интерфейс должен быть структурой или указателем на структуру",
			errors.ParamsOption("Тип интерфейса", reflectValue.Kind().String()),
			errors.SkipThisCallOption())
	}

	reflectType := reflectValue.Type()

	// Проходимся по каждому полю структуры
	for i := 0; i < reflectType.NumField(); i++ {

		// Получаем поле
		typeField := reflectType.Field(i)

		// Получаем значение поля
		valField := reflectValue.Field(i)

		// Получаем то, что в теге validate
		reqTag := typeField.Tag.Get("validate")

		// Получаем тег json
		jsTag := typeField.Tag.Get("json")

		// Если поле равно нулю и тег validate = required
		if valField.IsZero() && reqTag == "required" {

			if tag != "" {
				tag += "."
			}

			tags = append(tags, jsTag)
		}

		// Если тип поля структура
		if typeField.Type.Kind() == reflect.Struct {

			// Приводим к интерфейсу
			tt := valField.Interface()

			// Добавляем вложенность
			if len(tag) != 0 && i == 0 {
				tag += "."
			}

			// Рекурсивно вызываем функцию для вложенной функции
			_tags, err := zeroValue(ctx, tt, tag+jsTag, depth+1)

			// Если внутри структуры плохо
			if err != nil {
				return tags, err
			}

			tags = append(tags, _tags...)
		}
	}

	// Возвращаем гуд
	return tags, nil
}

func ZeroValue(ctx context.Context, requestStruct any) error {
	tags, err := zeroValue(ctx, requestStruct, "", stackTrace.Skip2PreviousCallers)
	if err != nil {
		return err
	}

	if tags != nil {
		params := make([]any, 0, len(tags)*2) //nolint:gomnd
		for _, tag := range tags {
			params = append(params, tag, "required")
		}
		return errors.BadRequest.New(ctx, "Required field is not filled",
			errors.ParamsOption(params...),
		)
	}

	return nil
}
