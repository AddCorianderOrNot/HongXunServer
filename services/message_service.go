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
	SaveMessage(message models.Message) models.Response
	LoadMessage(userTo, userFrom primitive.ObjectID) models.Response
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
		Keys:    bsonx.Doc{{"user_to", bsonx.Int32(1)}, {"user_from", bsonx.Int32(1)}},
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

func (s *messageService) SaveMessage(message models.Message) models.Response {
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
	return models.Response{
		ErrCode: 0,
		ErrMsg:  "发送消息成功",
		Data:    message,
	}
}

func (s *messageService) LoadMessage(userTo, userFrom primitive.ObjectID) models.Response {
	log.Println("LoadMessage", userTo, userFrom)
	var messages []*models.Message
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
		messages = append(messages, &elem)
	}
	err = cur.Close(context.TODO())
	return models.Response{
		ErrCode: 0,
		ErrMsg:  "加载消息成功",
		Data:    messages,
	}
}
