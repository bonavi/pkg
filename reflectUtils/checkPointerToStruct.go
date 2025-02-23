package reflectUtils

import (
	"context"
	"reflect"

	"pkg/errors"
)

func CheckPointerToStruct(ctx context.Context, dest any) error {

	// Проверяем типы данных
	reflectVar := reflect.ValueOf(dest)
	if reflectVar.Kind() != reflect.Ptr {
		return errors.InternalServer.New(ctx, "Пришедший интерфейс не является указателем",
			errors.ParamsOption("Тип интерфейса", reflectVar.Kind().String()),
			errors.SkipThisCallOption())
	}

	if reflectVar.Elem().Kind() != reflect.Struct {
		return errors.InternalServer.New(ctx, "Тип указателя не является структурой",
			errors.ParamsOption("Тип указателя", reflectVar.Kind().String()),
			errors.SkipThisCallOption())
	}

	return nil
}
