package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRoute(app *fiber.App) {
	app.Get("products", controllers.GetAllProducts)
	app.Post("product", controllers.CreateProduct)
	app.Get("product/:productId", controllers.GetDetailProduct)
}