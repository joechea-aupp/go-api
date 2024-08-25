package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		log.Println("error:", err)
		return err
	}

	if _, err := u.Collection.InsertOne(context.TODO(), User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}); err != nil {

		log.Println("error:", err)
		return err
	}

	return nil
}

func (u *UserService) GetUser(username string) (User, error) {
	var user User
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	err := u.Collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Println("error:", err)
		return User{}, err
	}

	return user, nil
}

func (u *UserService) GetUsers(start, limit int64) ([]User, error) {
	var users []User

	findOptions := options.Find().SetLimit(limit).SetSkip(start)

	cursor, err := u.Collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		panic(err)
	}

	return users, nil
}

func (u *UserService) TotalUsers() (int64, error) {
	count, err := u.Collection.EstimatedDocumentCount(context.TODO())
	if err != nil {
		log.Println("error:", err)
		return 0, err
	}

	return count, nil
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

func (u *UserService) DelUsers(ids []primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: ids}}}}

	_, err := u.Collection.DeleteMany(context.TODO(), filter)
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
