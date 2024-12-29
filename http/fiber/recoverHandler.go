package fiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"pkg/errors"
	"pkg/log"
)

func RecoverHandler(c *fiber.Ctx, e any) {

	err := errors.InternalServer.New(fmt.Sprintf("%v", e),
		errors.SkipPreviousCallerOption(),
	)

	if err = DefaultErrorHandler(c, err); err != nil {
		log.Error(c.Context(), err)
	}
}
