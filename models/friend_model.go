package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Friend struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	OwnerId    primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	FriendId   primitive.ObjectID `json:"friend_id" bson:"friend_id"`
	Remarks    string             `json:"remarks" bson:"remarks"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
	ReadTime   int64              `json:"read_time" bson:"read_time"`
}

type FriendEmail struct {
	Email string `json:"friend_id" url:"friend_id"`
}

type FriendTime struct {
	Email    string `json:"friend_id"`
	ReadTime int64  `json:"read_time"`
}
