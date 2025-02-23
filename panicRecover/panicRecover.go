package panicRecover

import (
	"context"
	"fmt"

	"pkg/errors"
)

func PanicRecover(ctx context.Context, handling func(err error)) {
	if r := recover(); r != nil {
		handling(errors.InternalServer.New(ctx, fmt.Sprintf("%v", r),
			errors.SkipThisCallOption(),
		))
	}
}
