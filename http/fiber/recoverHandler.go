package fiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"pkg/errors"
	"pkg/log"
	"pkg/middleware"
)

func RecoverHandler(c *fiber.Ctx, err any) {

	if err = middleware.DefaultErrorHandler(c, errors.Default.New(fmt.Sprintf("%v", err)).
		SkipPreviousCaller(),
	); err != nil {
		log.Error(err)
	}
}
