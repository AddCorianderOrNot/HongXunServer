package services

import (
	"HongXunServer/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

type UserService interface {
	Register(user models.User) error
}

type userService struct {
	C *mongo.Collection
}

func NewUserService(collection *mongo.Collection) UserService {
	log.Println("NewUserService")
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("userIndex").
		SetUnique(true).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{{
			Key:   "email",
			Value: bsonx.String("text"),
		}},
		Options: indexOpt,
	})
	if err != nil {
		log.Println(err)
	}

	return &userService{C: collection}
}

func (s *userService) Register(user models.User) error {
	log.Println("Register")
	err := s.C.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&models.User{})
	if err == nil {
		return err
	}
	if user.Id.IsZero() {
		user.Id = primitive.NewObjectID()
	}
	_, err = s.C.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
	}
	return err
}
