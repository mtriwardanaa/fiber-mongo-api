package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID primitive.ObjectID `json:"id,omitempty"`
	Name string `json:"name,omitempty" validate:"required"`
	Price int32 `json:"price,omitempty" validate:"required"`
}