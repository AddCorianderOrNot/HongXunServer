package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	CreateTime int                `json:"create_time" bson:"create_time"`
	UserFrom   string             `json:"user_from" bson:"user_from"`
	UserTo     string             `json:"user_to" bson:"user_to"`
	Content    string             `json:"content" bson:"content"`
	IsViewed   bool               `json:"is_viewed" bson:"is_viewed"`
}
