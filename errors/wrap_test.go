package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestWrapOverwritesType(t *testing.T) {

	// Типы для проверки, как в приложении
	notFound := ErrorType{Name: "NotFound", HTTPCode: http.StatusNotFound, LogAs: LogAsWarning, HumanText: "Не найдено"}
	tokenExpired := ErrorType{Name: "TokenExpired", HTTPCode: http.StatusUnauthorized, LogAs: LogAsWarning, HumanText: "Сессия устарела"}

	t.Run("1. Повторный Wrap перезаписывает тип и шаблонный текст", func(t *testing.T) {
		inner := Default.New("token is expired").WithParams("userID", "user-1")

		got := tokenExpired.Wrap(inner)

		if got.ErrorType != tokenExpired {
			t.Errorf("ErrorType: want %v, got %v", tokenExpired.Name, got.ErrorType.Name)
		}
		if got.HumanText != tokenExpired.HumanText {
			t.Errorf("HumanText: want %q, got %q", tokenExpired.HumanText, got.HumanText)
		}
		if got.Params["userID"] != "user-1" {
			t.Errorf("Params lost: %v", got.Params)
		}
	})

	t.Run("2. Кастомный текст для пользователя сохраняется", func(t *testing.T) {
		inner := notFound.New("no rows").WithCustomHumanText("Сервер не найден")

		got := tokenExpired.Wrap(inner)

		if got.ErrorType != tokenExpired || got.HumanText != "Сервер не найден" {
			t.Errorf("want type %v and custom text, got %v / %q", tokenExpired.Name, got.ErrorType.Name, got.HumanText)
		}
	})

	t.Run("3. Default.Wrap не затирает осмысленный тип", func(t *testing.T) {
		got := Default.Wrap(notFound.New("no rows"))

		if got.ErrorType != notFound {
			t.Errorf("want %v, got %v", notFound.Name, got.ErrorType.Name)
		}
	})

	t.Run("4. DontEraseErrorType фиксирует тип", func(t *testing.T) {
		got := tokenExpired.Wrap(notFound.New("no rows").DontEraseErrorType())

		if got.ErrorType != notFound {
			t.Errorf("want %v, got %v", notFound.Name, got.ErrorType.Name)
		}
	})

	t.Run("5. DontEraseErrorType после Wrap фиксирует тип, полученный этим Wrap", func(t *testing.T) {
		locked := tokenExpired.Wrap(notFound.New("no rows")).DontEraseErrorType()

		got := notFound.Wrap(locked)

		if got.ErrorType != tokenExpired {
			t.Errorf("want %v, got %v", tokenExpired.Name, got.ErrorType.Name)
		}
	})

	t.Run("6. Фиксация переживает опции и обёртку сторонней библиотеки через %w", func(t *testing.T) {
		locked := notFound.New("no rows").DontEraseErrorType().WithParams("id", "1").WithCustomHumanText("Не найдено")
		wrappedByLibrary := fmt.Errorf("library: %w", locked)

		got := tokenExpired.Wrap(wrappedByLibrary)

		if got.ErrorType != notFound {
			t.Errorf("want %v, got %v", notFound.Name, got.ErrorType.Name)
		}
	})

	t.Run("7. Необёрнутая ошибка оборачивается как раньше", func(t *testing.T) {
		got := notFound.Wrap(errors.New("plain"))

		if got.ErrorType != notFound || got.HumanText != notFound.HumanText {
			t.Errorf("want %v, got %v / %q", notFound.Name, got.ErrorType.Name, got.HumanText)
		}
	})
}
