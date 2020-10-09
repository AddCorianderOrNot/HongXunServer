package services

import (
	"HongXunServer/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MessageService interface {
	SaveMessage(message models.Message) error
	LoadMessage() ([]*models.Message, error)
}

type messageService struct {
	C *mongo.Collection
}

func NewMessageService(collection *mongo.Collection) MessageService {
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("messageIndex").
		SetUnique(true).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    []string{"Id"},
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
	log.Println("SaveMessage")
	insertResult, err := s.C.InsertOne(context.TODO(), message)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Inserted a single message: ", insertResult.InsertedID)
	return err
}

func (s *messageService) LoadMessage() ([]*models.Message, error) {
	log.Println("LoadMessage")
	var Messages []*models.Message
	findOptions := options.Find().SetLimit(10)
	cur, err := s.C.Find(context.TODO(), bson.D{{"isviewed", false}}, findOptions)
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
