package models

import (
	"time"
	"github.com/appu900/authservice/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	LastLogin *time.Time         `json:"last_login,omitempty" bson:"last_login,omitempty"`
	Active    bool               `json:"active" bson:"active" default:"true"`
}

func NewUser(username, password, email string) *User {
	now := time.Now()
	return &User{
		ID:        primitive.NewObjectID(),
		Username:  username,
		Password:  password,
		CreatedAt: now,
		Email:     email,
		UpdatedAt: now,
		Active:    true,
	}
}

func (u *User) BeforeInsert() {
	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

func CheckIfUserExistsWithEmail(c *fiber.Ctx, email string) (*User, error) {
	userCollection := database.GetCollection("users")
	var existingUser User
	err := userCollection.FindOne(c.Context(), bson.M{
		"email": email,
	}).Decode(&existingUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &existingUser, err
}
