package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Nickname  string             `json:"nickname" bson:"nickname"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Icon      string             `json:"icon" bson:"icon"`
	Token     string             `json:"token" bson:"token"`
	Signature string             `json:"signature" bson:"signature"`
}

type UserMini struct {
	Nickname  string `json:"nickname" bson:"nickname"`
	Email     string `json:"email" bson:"email"`
	Icon      string `json:"icon" bson:"icon"`
	Signature string `json:"signature" bson:"signature"`
}

type Authentication struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
