package jwtManager

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"pkg/errors"
)

var (
	ErrTokenExpired = errors.New("tokenExpired")
)

type JWTManager struct {
	accessTokenSigningKey []byte
	ttls                  map[TokenType]time.Duration
}

var jwtManager *JWTManager

func Init(
	accessTokenSigningKey []byte,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) {
	jwtManager = &JWTManager{
		accessTokenSigningKey: accessTokenSigningKey,
		ttls: map[TokenType]time.Duration{
			RefreshToken: refreshTokenTTL,
			AccessToken:  accessTokenTTL,
		},
	}
}

type MyCustomClaims[T any] struct {
	CustomClaims T
	jwt.StandardClaims
}

type TokenType int

const (
	RefreshToken = iota + 1
	AccessToken
)

func GenerateToken[T any](ctx context.Context, tokenType TokenType, customClaims T) (string, error) {

	if jwtManager == nil {
		return "", errors.InternalServer.New(ctx, "JWTManager is not initialized")
	}

	claims := MyCustomClaims[T]{
		CustomClaims: customClaims,
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: time.Now().Add(jwtManager.ttls[tokenType]).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bonavii.com",
			NotBefore: 0,
			Subject:   "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtManager.accessTokenSigningKey)
	if err != nil {
		return "", errors.InternalServer.Wrap(ctx, err)
	}

	return tokenStr, nil
}

func ParseToken[T any](ctx context.Context, reqToken string) (T, error) {

	var typeZeroValue T

	// Если менеджер не инициализирован, возвращаем ошибку
	if jwtManager == nil {
		return typeZeroValue, errors.InternalServer.New(ctx, "JWTManager is not initialized",
			errors.SkipThisCallOption(),
		)
	}

	// Если токен пустой, возвращаем ошибку
	if reqToken == "" {
		return typeZeroValue, errors.Unauthorized.New(ctx, "JWT-token is empty",
			errors.SkipThisCallOption(),
		)
	}

	// Парсим токен
	token, jwtErr := jwt.ParseWithClaims(reqToken, &MyCustomClaims[T]{}, func(token *jwt.Token) (i any, err error) { //nolint:exhaustruct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InternalServer.New(ctx, "Unexpected signing method",
				errors.ParamsOption("token", token.Header["alg"]),
				errors.SkipThisCallOption(),
			)
		}

		return jwtManager.accessTokenSigningKey, nil
	})

	// Если ошибка не пустая
	if jwtErr != nil {

		// Если ошибка валидатора
		var validationError *jwt.ValidationError
		if errors.As(jwtErr, &validationError) {

			switch {

			// Если токен истек
			case validationError.Errors == jwt.ValidationErrorExpired:

				// Если токен истек, определяем ошибку с errorf, чтобы потом вернуть
				jwtErr = errors.Unauthorized.Wrap(ctx, jwtErr,
					errors.SkipPreviousCallerOption(),
					errors.ErrorfOption(ErrTokenExpired),
				)

			// Если другая ошибка, просто возвращаем ее
			default:
				return typeZeroValue, errors.Unauthorized.Wrap(ctx, jwtErr,
					errors.SkipPreviousCallerOption(),
				)
			}

		} else { // Если ошибка не валидатора
			return typeZeroValue, errors.InternalServer.Wrap(ctx, jwtErr,
				errors.SkipPreviousCallerOption(),
			)
		}
	}

	// Если ошибок нет, пробуем получить кастомные клеймы
	claims, ok := token.Claims.(*MyCustomClaims[T])
	if !ok {
		return typeZeroValue, errors.InternalServer.New(ctx, "Error get user claims from token")
	}

	// Обрабатываем ошибку парсера jwt
	if jwtErr != nil {
		return claims.CustomClaims, jwtErr
	}

	return claims.CustomClaims, nil
}
