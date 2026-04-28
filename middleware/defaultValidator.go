package middleware

import (
	"pkg/errors"
	"pkg/validator"
)

type validatorProtocol interface {
	Validate() error
}

func DefaultValidator(object any) (err error) {

	// Если структура реализует интерфейс валидатора, то валидируем ее с помощью функции
	if v, ok := object.(validatorProtocol); ok {
		if err = v.Validate(); err != nil {
			return errors.Default.Wrap(err).SkipThisCall()
		}
	}

	// Валидируем структуру с помощью декларативного валидатора по тегам
	if err = validator.Validate(object); err != nil {
		return errors.Default.Wrap(err).SkipThisCall()
	}

	return nil
}
