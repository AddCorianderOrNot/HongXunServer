package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/base64"
)

func Md5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64(str []byte) string {
	return base64.StdEncoding.EncodeToString(str)
}