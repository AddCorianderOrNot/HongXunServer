package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	CreateTime time.Time          `json:"dateTime" bson:"create_time"`
	UserFrom   string 			  `json:"sender" bson:"user_from"`
	UserName   string 			  `json:"senderName" bson:"user_from"`
	UserTo     string             `json:"receiver" bson:"user_to"`
	Type       int                `json:"messageType" bson:"type"`
	Content    string             `json:"message" bson:"content"`
	IsViewed   bool               `json:"is_viewed" bson:"is_viewed"`
}
