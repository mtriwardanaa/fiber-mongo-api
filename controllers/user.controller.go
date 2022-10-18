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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validateUser = validator.New()

func CreateUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var user models.User
    defer cancel()
  
    //validateUser the request body
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    //use the validator library to validateUser required fields
    if validationErr := validateUser.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"result": validationErr.Error()}})
    }
  
    newUser := models.User{
        Id:       primitive.NewObjectID(),
        Name:     user.Name,
        Location: user.Location,
        Title:    user.Title,
    }
  
    result, err := userCollection.InsertOne(ctx, newUser)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"result": result}})
}

func GetAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    var user models.User
    defer cancel()
  
    objId, _ := primitive.ObjectIDFromHex(userId)
  
    err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"result": user}})
}

func EditAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    var user models.User
    defer cancel()
  
    objId, _ := primitive.ObjectIDFromHex(userId)
  
    //validateUser the request body
    if err := c.BodyParser(&user); err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    //use the validator library to validateUser required fields
    if validationErr := validateUser.Struct(&user); validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"result": validationErr.Error()}})
    }
  
    update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
  
    result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    //get updated user details
    var updatedUser models.User
    if result.MatchedCount == 1 {
        err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
        }
    }
  
    return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"result": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    userId := c.Params("userId")
    defer cancel()
  
    objId, _ := primitive.ObjectIDFromHex(userId)
  
    result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    if result.DeletedCount < 1 {
        return c.Status(http.StatusNotFound).JSON(
            responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"result": "User with specified ID not found!"}},
        )
    }
  
    return c.Status(http.StatusOK).JSON(
        responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"result": "User successfully deleted!"}},
    )
}

func GetAllUsers(c *fiber.Ctx) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    var users []models.User
    defer cancel()
  
    results, err := userCollection.Find(ctx, bson.M{})
  
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
    }
  
    //reading from the db in an optimal way
    defer results.Close(ctx)
    for results.Next(ctx) {
        var singleUser models.User
        if err = results.Decode(&singleUser); err != nil {
            return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"result": err.Error()}})
        }
      
        users = append(users, singleUser)
    }
  
    return c.Status(http.StatusOK).JSON(
        responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"result": users}},
    )
}