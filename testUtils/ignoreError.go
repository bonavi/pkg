package testUtils

import (
	"pkg/errors"
	"pkg/log"
)

func IgnoreErrorWithArgument[T any](v T, err error) T {
	if err != nil {
		log.Fatal(errors.InternalServer.Wrap(err).SkipThisCall())
	}
	return v
}

func IgnoreError(err error) {
	if err != nil {
		log.Fatal(errors.InternalServer.Wrap(err).SkipThisCall())
	}
}
