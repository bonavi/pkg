package fiber

import "github.com/gofiber/fiber/v2"

func NewReadyHandler(ready chan struct{}) func(*fiber.Ctx) bool {
	return func(c *fiber.Ctx) bool {
		select {
		case <-ready:
			return true
		default:
			return false
		}
	}
}
