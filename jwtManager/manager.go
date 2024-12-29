package jwtManager

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"pkg/errors"
)

var (
	ErrUserUnauthorized = errors.New("user unauthorized")
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

func GenerateToken[T any](tokenType TokenType, customClaims T) (string, error) {

	if jwtManager == nil {
		return "", errors.InternalServer.New("JWTManager is not initialized")
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
		return "", errors.InternalServer.Wrap(err)
	}

	return tokenStr, nil
}

func ParseToken[T any](reqToken string) (T, error) {

	var typeZeroValue T

	if jwtManager == nil {
		return typeZeroValue, errors.InternalServer.New("JWTManager is not initialized",
			errors.SkipThisCallOption(),
		)
	}

	if reqToken == "" {
		return typeZeroValue, errors.Unauthorized.New("JWT-token is empty",
			errors.SkipThisCallOption(),
		)
	}

	token, jwtErr := jwt.ParseWithClaims(reqToken, &MyCustomClaims[T]{}, func(token *jwt.Token) (i any, err error) { //nolint:exhaustruct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InternalServer.New("Unexpected signing method",
				errors.ParamsOption("token", token.Header["alg"]),
				errors.SkipThisCallOption(),
			)
		}

		return jwtManager.accessTokenSigningKey, nil
	})
	if jwtErr != nil {
		var validationError *jwt.ValidationError
		if errors.As(jwtErr, &validationError) {

			switch {
			case validationError.Errors == jwt.ValidationErrorExpired:
				jwtErr = errors.Unauthorized.Wrap(jwtErr,
					errors.SkipPreviousCallerOption(),
					errors.ErrorfOption(ErrUserUnauthorized),
				)
			default:
				return typeZeroValue, errors.Unauthorized.Wrap(jwtErr,
					errors.SkipPreviousCallerOption(),
				)
			}

		} else {
			return typeZeroValue, errors.InternalServer.Wrap(jwtErr,
				errors.SkipPreviousCallerOption(),
			)
		}
	}

	claims, ok := token.Claims.(*MyCustomClaims[T])
	if !ok {
		return typeZeroValue, errors.InternalServer.New("Error get user claims from token")
	}

	return claims.CustomClaims, jwtErr
}
