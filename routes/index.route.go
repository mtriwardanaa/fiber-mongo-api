package routes

import "github.com/gofiber/fiber/v2"

func IndexRoute(app *fiber.App) {
	UserRoute(app)
	ProductRoute(app)
}