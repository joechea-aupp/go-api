package factory

import (
	"context"
	"log"
	"time"

	"github.com/joechea-aupp/go-api/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// when declare method to a pointer, it is call pointer receiver.
// pointer receiver is used to modify the value of the receiver.
func (u *User) CreateUser(user User) error {
	collection := service.Collection("users")
	_, err := collection.InsertOne(context.TODO(), User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error:", err)
		return err
	}

	return nil
}

func (u *User) GetUser(email string) (User, error) {
	var user User

	collection := service.Collection("users")
	filter := bson.D{primitive.E{Key: "email", Value: email}}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println("error:", err)
		return User{}, err
	}

	return user, nil
}
