package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UserService struct {
	Collection *mongo.Collection
}

func NewUserService() *UserService {
	return &UserService{
		Collection: collection("users"),
	}
}

// when declare method to a pointer, it is call pointer receiver.
// pointer receiver is used to modify the value of the receiver.
func (u *UserService) CreateUser(user User) error {
	_, err := u.Collection.InsertOne(context.TODO(), User{
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

func (u *UserService) GetUser(email string) (User, error) {
	var user User
	filter := bson.D{primitive.E{Key: "email", Value: email}}

	err := u.Collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println("error:", err)
		return User{}, err
	}

	return user, nil
}

func (u *UserService) GetUsers() ([]User, error) {
	var users []User

	cursor, err := u.Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		panic(err)
	}

	return users, nil
}

func (u *UserService) DelUser(id string) error {
	// convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	_, err = u.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) UpdateUser(id string, data User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectID}}

	update := bson.D{
		{Key: "$set", Value: data},
	}

	if _, err := u.Collection.UpdateOne(context.TODO(), filter, update); err != nil {
		return err
	}

	return nil
}
