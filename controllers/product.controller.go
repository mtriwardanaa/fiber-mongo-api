package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")
var validateProduct = validator.New()

func CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var product models.Product
	defer cancel()

	if err := c.BodyParser(&product); err != nil {
		return responses.ErrorResponse(c, err)
	}
}

func GetAllProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var products []models.Product
	defer cancel()

	results, err := productCollection.Find(ctx, bson.M{})

	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				return responses.ErrorResponse(c, err)
			}
		
			products = append(products, singleProduct)
	}

	return c.Status(http.StatusOK).JSON(
			responses.GenResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"result": products}},
	)
}

func GetDetailProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	var product models.Product
	var productId = c.Params("productId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(productId)

	err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&product)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"result": product}})
}