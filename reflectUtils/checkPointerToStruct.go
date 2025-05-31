package reflectUtils

import (
	"reflect"

	"pkg/errors"
)

func CheckPointerToStruct(dest any) error {

	// Проверяем типы данных
	reflectVar := reflect.ValueOf(dest)
	if reflectVar.Kind() != reflect.Ptr {
		return errors.InternalServer.New("Пришедший интерфейс не является указателем").
			WithParams("Тип интерфейса", reflectVar.Kind().String()).
			WithStackTraceJump(errors.SkipThisCall)
	}

	if reflectVar.Elem().Kind() != reflect.Struct {
		return errors.InternalServer.New("Тип указателя не является структурой").
			WithParams("Тип указателя", reflectVar.Kind().String()).
			WithStackTraceJump(errors.SkipThisCall)
	}

	return nil
}
