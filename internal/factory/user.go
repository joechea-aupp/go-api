package factory

import (
	"context"
	"log"
	"time"

	"github.com/joechea-aupp/go-api/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	CreatedAt time.Time          `bson:"created_at"`
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
