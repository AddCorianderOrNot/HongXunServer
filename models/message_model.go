package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	CreateTime int64              `json:"dateTime" bson:"create_time"`
	UserFrom   string             `json:"sender" bson:"user_from"`
	UserName   string             `json:"senderName" bson:"UserName"`
	UserTo     string             `json:"receiver" bson:"user_to"`
	Type       int                `json:"messageType" bson:"type"`
	Content    string             `json:"message" bson:"content"`
	IsViewed   bool               `json:"is_viewed" bson:"is_viewed"`
}
