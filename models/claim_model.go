package models

import "github.com/kataras/iris/v12/middleware/jwt"

type UserClaims struct {
	jwt.Claims
	Email string
}
