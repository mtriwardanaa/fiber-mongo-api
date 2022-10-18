package responses

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusInternalServerError).JSON(GenResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
}