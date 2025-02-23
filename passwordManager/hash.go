package passwordManager

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"pkg/errors"
)

const (
	saltSize          = 16
	maxPasswordLength = 72
)

func CreateNewPassword(ctx context.Context, password, generalSalt, userSalt []byte) ([]byte, error) {

	if len(password) > maxPasswordLength {
		availableLength := len(password) - len(generalSalt) - saltSize
		return nil, errors.BadRequest.New(ctx, fmt.Sprintf("Длина пароля не должна превышать %v символов", availableLength))
	}

	passwordHash, err := bcrypt.GenerateFromPassword(SaltPassword(password, userSalt, generalSalt), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	return passwordHash, nil
}

func CompareHashAndPassword(ctx context.Context, hash, password, userSalt, generalSalt []byte) error {

	if len(password) > maxPasswordLength {
		availableLength := len(password) - len(generalSalt) - saltSize
		return errors.BadRequest.New(ctx, fmt.Sprintf("Длина пароля не должна превышать %v символов", availableLength))
	}

	if err := bcrypt.CompareHashAndPassword(hash, SaltPassword(password, userSalt, generalSalt)); err != nil {
		return errors.BadRequest.Wrap(ctx, err,
			errors.HumanTextOption("Неверно введен логин или пароль"),
		)
	}
	return nil
}

func SaltPassword(password, userSalt, generalSalt []byte) []byte {
	return bytes.Join([][]byte{userSalt, password, generalSalt}, nil)
}

func GenerateRandomSalt(ctx context.Context) ([]byte, error) {
	b := make([]byte, saltSize)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.InternalServer.Wrap(ctx, err)
	}

	return b, nil
}
