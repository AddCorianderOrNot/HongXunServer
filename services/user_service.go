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
	Register(user models.User) bool
	Verify(authentication models.Authentication) bool
	isExist(email string) (bool, models.User)
}

type userService struct {
	C *mongo.Collection
}

func NewUserService(collection *mongo.Collection) UserService {
	log.Println("NewUserService")
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("userIndex").
		SetUnique(false).
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

func (s *userService) isExist(email string) (bool, models.User) {
	log.Println("Find:", email)
	var user models.User
	err := s.C.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	log.Println(user, err)
	if err == nil {
		return true, user
	} else {
		return false, user
	}
}

func (s *userService) Register(user models.User) bool {
	log.Println("Register")
	exist, _ := s.isExist(user.Email)
	log.Println(exist)
	if exist {
		return false
	}
	if user.Id.IsZero() {
		user.Id = primitive.NewObjectID()
	}
	log.Println("Insert:", user)
	_, err := s.C.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
	}
	return true
}

func (s *userService) Verify(authentication models.Authentication) bool {
	exist, user := s.isExist(authentication.Email)
	log.Println(exist, user)
	if exist && authentication.Password == user.Password {
		return true
	} else {
		return false
	}
}
