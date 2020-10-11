package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
	UserFrom   primitive.ObjectID `json:"user_from" bson:"user_from"`
	UserTo     primitive.ObjectID `json:"user_to" bson:"user_to"`
	Type       int                `json:"type" bson:"type"`
	Content    string             `json:"content" bson:"content"`
	IsViewed   bool               `json:"is_viewed" bson:"is_viewed"`
}
