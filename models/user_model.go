package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Nickname  string             `json:"nickname"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	Icon      string             `json:"icon"`
	Token     string             `json:"token"`
	Signature string             `json:"signature"`
}

type UserMini struct {
	Nickname  string             `json:"nickname"`
	Email     string             `json:"email"`
	Icon      string             `json:"icon"`
	Signature string             `json:"signature"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
