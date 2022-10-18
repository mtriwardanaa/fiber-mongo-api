package responses

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func OKResponse(c *fiber.Ctx, data error) error {
	return c.Status(http.StatusInternalServerError).JSON(GenResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
}