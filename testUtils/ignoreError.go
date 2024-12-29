package testUtils

import (
	"context"

	"pkg/errors"
	"pkg/log"
)

func IgnoreErrorWithArgument[T any](v T, err error) T {
	if err != nil {
		log.Fatal(context.Background(), errors.InternalServer.Wrap(err, errors.SkipThisCallOption()))
	}
	return v
}

func IgnoreError(err error) {
	if err != nil {
		log.Fatal(context.Background(), errors.InternalServer.Wrap(err, errors.SkipThisCallOption()))
	}
}
