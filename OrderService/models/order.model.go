package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RestaurantID string             `json:"restaurant_id,omitempty" bson:"restaurant_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	
}
