package middleware

import (
	"HongXunServer/auth"
	"HongXunServer/models"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

var J = jwt.HMAC(6000*time.Minute, auth.Secret, auth.Itsa16bytesecret)
var Verify = J.Verify


func GenerateToken(nickname string, id primitive.ObjectID, email string) string {
	log.Println("密码正确")
	claims := models.UserClaims{
		Claims: J.Expiry(jwt.Claims{
			Issuer:   "HongXun",
			Audience: jwt.Audience{nickname},
		}),
		UserId: id,
		UserEmail: email,
	}
	accessToken, _ := J.Token(claims)
	return accessToken
}