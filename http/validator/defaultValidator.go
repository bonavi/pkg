package middleware

import (
	"context"

	"pkg/errors"
	"pkg/validator"
)

type validatorProtocol interface {
	Validate(ctx context.Context) error
}

func Validate(ctx context.Context, object any) (err error) {

	// Если структура реализует интерфейс валидатора, то валидируем ее с помощью функции
	if v, ok := object.(validatorProtocol); ok {
		if err = v.Validate(ctx); err != nil {
			return errors.BadRequest.Wrap(ctx,
				err,
				errors.SkipThisCallOption(),
			)
		}
	}

	// Валидируем структуру с помощью декларативного валидатора по тегам
	if err = validator.Validate(object); err != nil {
		return errors.BadRequest.Wrap(ctx, err,
			errors.SkipThisCallOption(),
		)
	}

	return nil
}
