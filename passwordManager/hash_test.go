package passwordManager

import (
	"fmt"
	"testing"

	"pkg/errors"
	"pkg/testUtils"
)

func TestCompareHashAndPassword(t *testing.T) {
	type args struct {
		hash        []byte
		password    []byte
		userSalt    []byte
		generalSalt []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompareHashAndPassword(tt.args.hash, tt.args.password, tt.args.userSalt, tt.args.generalSalt); (err != nil) != tt.wantErr {
				t.Errorf("CompareHashAndPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateNewPassword(t *testing.T) {
	generalUserSalt := []byte("1234")
	t.Run("1. Корректное создание пароля", func(t *testing.T) {

		userPassword := []byte("password")
		generalSalt := []byte("generalSalt")

		// Генерируем новый пароль
		passwordHash, gotErr := CreateNewPassword(userPassword, generalSalt, generalUserSalt)
		testUtils.CheckError(t, nil, gotErr, false)

		// Проверяем, что пароль сгенерирован корректно
		gotErr = CompareHashAndPassword(passwordHash, userPassword, generalUserSalt, generalSalt)
		testUtils.CheckError(t, nil, gotErr, false)
	})
	t.Run("2. Ошибка генерации пароля из-за длины больше 72 символов", func(t *testing.T) {

		userPassword := []byte("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
		generalSalt := []byte("generalSalt")

		// Генерируем новый пароль
		_, gotErr := CreateNewPassword(userPassword, generalSalt, generalUserSalt)
		testUtils.CheckError(t, errors.BadRequest.New(""), gotErr, false)
	})
}

func TestGenerateRandomBytes(t *testing.T) {

	count := 1000
	mapOfHashes := make(map[string]struct{}, count)

	for i := 1; i <= count; i++ {
		t.Run(fmt.Sprintf("%v. Создание уникального токена", i), func(t *testing.T) {
			hash, err := GenerateRandomSalt()
			if testUtils.CheckError(t, nil, err, false) {
				return
			}
			if _, ok := mapOfHashes[string(hash)]; ok {
				t.Errorf("Сгенерированный хэш не уникален. Хэш: %s", hash)
			}
			mapOfHashes[string(hash)] = struct{}{}
		})
	}
}
