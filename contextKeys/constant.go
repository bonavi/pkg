package contextKeys

import (
	"context"
)

type ContextKey int

const (
	DeviceIDKey ContextKey = iota + 1
	UserIDKey
	TaskIDKey
	XRequestIDKey
)

func SetRequestID(ctx context.Context, taskID string) context.Context {
	return context.WithValue(ctx, TaskIDKey, taskID)
}

func GetRequestID(ctx context.Context) *string {
	if v, ok := ctx.Value(TaskIDKey).(string); ok {
		return &v
	}
	return nil
}
