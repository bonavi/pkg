package env

import (
	"sync/atomic"

	"github.com/caarlos0/env/v11"

	"pkg/errors"
	"pkg/log"
)

var instance any
var alreadyCalled atomic.Bool

// Load returns a new ConfigType.
func Load[T any]() T {

	// Если уже вызывали, то возвращаем тот же инстанс
	if alreadyCalled.Load() {
		conf, ok := instance.(*T)
		if !ok {
			log.Fatal(errors.InternalServer.New("Config isn't type of T"))
		}
		return *conf
	}

	// Обозначаем, что уже вызывали эту функцию
	alreadyCalled.Store(true)

	// Создаем новый инстанс
	instance = new(T)

	// Парсим env'ы в инстанс
	if err := env.ParseWithOptions(instance, env.Options{
		Environment:           nil,
		TagName:               "",
		PrefixTagName:         "",
		DefaultValueTagName:   "",
		RequiredIfNoDef:       true,
		OnSet:                 nil,
		Prefix:                "",
		UseFieldNameByDefault: false,
		FuncMap:               nil,
	}); err != nil {
		log.Fatal(errors.InternalServer.Wrap(err).WithStackTraceJump(errors.SkipPreviousCaller))
	}

	conf, ok := instance.(*T)
	if !ok {
		log.Fatal(errors.InternalServer.New("Config isn't type of T"))
	}
	return *conf
}

// LoadForTest returns a new ConfigType.
func LoadForTest[T any]() T {

	// Если уже вызывали, то возвращаем тот же инстанс
	if alreadyCalled.Load() {
		conf, ok := instance.(*T)
		if !ok {
			log.Fatal(errors.InternalServer.New("Config isn't type of T"))
		}
		return *conf
	}

	// Обозначаем, что уже вызывали эту функцию
	alreadyCalled.Store(true)

	// Создаем новый инстанс
	instance = new(T)

	// Парсим env'ы в инстанс
	if err := env.ParseWithOptions(instance, env.Options{
		Environment:           nil,
		TagName:               "",
		PrefixTagName:         "",
		DefaultValueTagName:   "",
		RequiredIfNoDef:       false,
		OnSet:                 nil,
		Prefix:                "",
		UseFieldNameByDefault: false,
		FuncMap:               nil,
	}); err != nil {
		log.Fatal(errors.InternalServer.Wrap(err).WithStackTraceJump(errors.SkipPreviousCaller))
	}

	conf, ok := instance.(*T)
	if !ok {
		log.Fatal(errors.InternalServer.New("Config isn't type of T"))
	}
	return *conf
}

func Reset() {

	// Сбрасываем положение индикатора вызова функции
	alreadyCalled.Store(false)
}
