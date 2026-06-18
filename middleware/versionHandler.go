package middleware

import (
	"github.com/gofiber/fiber/v2"

	"pkg/errors"
)

type versionRes struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	Hostname  string `json:"hostname"`
	BuildDate string `json:"buildDate"`
}

func NewVersionHandler(version string, build string, buildDate string, hostname string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		if err := ctx.Status(fiber.StatusOK).JSON(versionRes{
			Version:   version,
			Build:     build,
			Hostname:  hostname,
			BuildDate: buildDate,
		}); err != nil {
			return errors.Default.Wrap(err)
		}
		return nil
	}
}
