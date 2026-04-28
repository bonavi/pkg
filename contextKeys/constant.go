package contextKeys

import (
	"context"

	"pkg/errors"
)

type ContextKey int

const (
	XRequestIDKey ContextKey = iota + 1
)

func GetXRequestID(ctx context.Context) (string, error) {
	xRequestID, ok := ctx.Value(XRequestIDKey).(string)
	if !ok {
		return "", errors.Default.New("XRequestID was not found").SkipThisCall()
	}
	return xRequestID, nil
}
