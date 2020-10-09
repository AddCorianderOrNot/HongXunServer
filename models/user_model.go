package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	Icon      string             `json:"icon"`
	Signature string             `json:"signature"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
