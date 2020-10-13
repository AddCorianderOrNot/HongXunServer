package middleware

import (
	"HongXunServer/auth"
	"github.com/kataras/iris/v12/middleware/jwt"
	"time"
)

var J = jwt.HMAC(15*time.Minute, auth.Secret, auth.Itsa16bytesecret)
var Verify = J.Verify
