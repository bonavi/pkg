package authManager

import (
	"context"
	"pkg/errors"
)

type AuthMethod interface {
	GetToken(ctx context.Context, id string) (string, error)
}

type AuthManager struct {
	authMethod AuthMethod

	accessTokensToServers map[string]string
}

func NewAuthManager(
	authMethod AuthMethod,
) *AuthManager {
	return &AuthManager{
		authMethod:            authMethod,
		accessTokensToServers: make(map[string]string),
	}
}

func (a *AuthManager) GetToken(ctx context.Context, id string) (string, error) {

	// Получаем токен доступа для сервера по его ID
	token, ok := a.accessTokensToServers[id]

	// Если токен не найден, вызываем метод авторизации
	if !ok {

		// Проверяем, установлен ли authMethod
		if a.authMethod == nil {
			return "", errors.Default.New("authMethod is not set")
		}

		// Делаем запрос к authMethod для получения токена
		var err error
		token, err = a.authMethod.GetToken(ctx, id)
		if err != nil {
			return "", err
		}

		// Сохраняем токен в мапу для последующего использования
		a.accessTokensToServers[id] = token
	}

	return token, nil
}
