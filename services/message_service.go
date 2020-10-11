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
	"time"
)

type MessageService interface {
	SaveMessage(message models.Message) error
	LoadMessage(userTo, userFrom primitive.ObjectID) ([]*models.Message, error)
}

type messageService struct {
	C *mongo.Collection
}

func NewMessageService(collection *mongo.Collection) MessageService {
	log.Println("NewMessageService")
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("messageIndex").
		SetUnique(false).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{{
			Key:   "user_to",
			Value: bsonx.String("text"),
		}},
		Options: indexOpt,
	})
	if err != nil {
		log.Println(err)
	}

	return &messageService{C: collection}
}

//func FindAll() []models.Message {
//	collection.Find()
//}

func (s *messageService) SaveMessage(message models.Message) error {
	if message.Id.IsZero() {
		message.Id = primitive.NewObjectID()
	}
	message.CreateTime = time.Now()
	message.IsViewed = false
	log.Println("SaveMessage", message)
	_, err := s.C.InsertOne(context.TODO(), message)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *messageService) LoadMessage(userTo, userFrom primitive.ObjectID) ([]*models.Message, error) {
	log.Println("LoadMessage", userTo, userFrom)
	var Messages []*models.Message
	findOptions := options.Find().SetLimit(10)
	cur, err := s.C.Find(context.TODO(), bson.D{
		{"is_viewed", false},
		{"user_to", userTo},
		{"user_from", userFrom},
	}, findOptions)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.Message
		_ = cur.Decode(&elem)
		Messages = append(Messages, &elem)
	}
	cur.Close(context.TODO())
	return Messages, err
}
