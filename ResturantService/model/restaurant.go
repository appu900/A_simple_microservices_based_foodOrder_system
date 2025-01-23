package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeoJSON struct {
	Type        string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}

type Dish struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Photo       string             `json:"photo,omitempty" bson:"photo,omitempty"`
	Avaliable   bool               `json:"avaliable,omitempty" bson:"avaliable,omitempty"`
}

type Restaurant struct {
	ID         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name,omitempty" bson:"name,omitempty"`
	Menu       []primitive.ObjectID `json:"menu,omitempty" bson:"menu,omitempty"` // Dish IDs
	Photo      string               `json:"photo,omitempty" bson:"photo,omitempty"`
	Location   GeoJSON              `json:"location,omitempty" bson:"location,omitempty"`
	Address    string               `json:"address,omitempty" bson:"address,omitempty"`
	Created_At time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_At time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

func NewDish(name string, price float64, description string, photoUrl string) *Dish {
	return &Dish{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Price:       price,
		Description: description,
		Photo:       photoUrl,
	}
}

func NewRestaurant(name, photoUrl, address string, longitude, lattitude float64) *Restaurant {
	now := time.Now()
	return &Restaurant{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Photo:      photoUrl,
		Address:    address,
		Created_At: now,
		Updated_At: now,
		Menu:       make([]primitive.ObjectID, 0),
		Location: GeoJSON{
			Type:        "Point",
			Coordinates: []float64{longitude, lattitude}, // Default location [0, 0]
		},
	}
}
