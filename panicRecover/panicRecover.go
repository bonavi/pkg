package panicRecover

import (
	"fmt"
	"pkg/errors"
)

func PanicRecover(handling func(err error)) {
	if r := recover(); r != nil {
		handling(errors.Default.New(fmt.Sprintf("%v", r)).
			SkipThisCall(),
		)
	}
}
