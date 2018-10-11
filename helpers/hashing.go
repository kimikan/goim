package helpers

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
)

func MD5(buf []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(buf)
	cipherStr := md5Ctx.Sum(nil)

	return cipherStr
	//return hex.EncodeToString(cipherStr)
}

func Base64Encode(src []byte) string {
	ret := base64.StdEncoding.EncodeToString(src)
	return strings.Replace(ret, "/", "-", -1)
}
