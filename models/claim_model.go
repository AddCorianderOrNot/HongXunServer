package models

import (
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	jwt.Claims
	UserId primitive.ObjectID
	UserEmail string
}
