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

const (
	addSuccessCode = 0
	addSuccessMsg  = "添加成功"
	addSelfCode    = 1
	addSelfMsg     = "不能添加自己为好友"
	addErrorCode   = 2
	addErrorMsg    = "未知错误"
)

type FriendService interface {
	AddFriend(ownerId, friendId primitive.ObjectID) models.Response
	LoadFriend(ownerId primitive.ObjectID) models.Response
}

type friendService struct {
	C *mongo.Collection
}

func NewFriendService(collection *mongo.Collection) FriendService {
	log.Println("NewFriendService")
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("friendIndex").
		SetUnique(false).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bsonx.Doc{{"owner_id", bsonx.Int32(1)}, {"friend_id", bsonx.Int32(1)}},
		Options: indexOpt,
	})
	if err != nil {
		log.Println(err)
	}

	return &friendService{C: collection}
}

func (s *friendService) AddFriend(ownerId, friendId primitive.ObjectID) models.Response {
	if ownerId == friendId {
		return models.Response{
			ErrCode: addSelfCode,
			ErrMsg:  addSelfMsg,
			Data:    nil,
		}
	}
	_, err := s.C.InsertOne(context.TODO(), models.Friend{
		Id:         primitive.NewObjectID(),
		OwnerId:    ownerId,
		FriendId:   friendId,
		CreateTime: time.Now(),
	})

	_, err = s.C.InsertOne(context.TODO(), models.Friend{
		Id:         primitive.NewObjectID(),
		OwnerId:    friendId,
		FriendId:   ownerId,
		CreateTime: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return models.Response{
			ErrCode: addErrorCode,
			ErrMsg:  addErrorMsg,
			Data:    nil,
		}
	}
	return models.Response{
		ErrCode: addSuccessCode,
		ErrMsg:  addSuccessMsg,
		Data:    nil,
	}
}

func (s *friendService) LoadFriend(ownerId primitive.ObjectID) models.Response {
	var friends []*models.Friend
	findOptions := options.Find().SetLimit(10)
	cur, err := s.C.Find(context.TODO(), bson.D{
		{"owner_id", ownerId},
	}, findOptions)
	if err != nil {
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var elem models.Friend
		_ = cur.Decode(&elem)
		friends = append(friends, &elem)
	}
	err = cur.Close(context.TODO())
	return models.Response{
		ErrCode: 0,
		ErrMsg:  "成功",
		Data:    friends,
	}
}
