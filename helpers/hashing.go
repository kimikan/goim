package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(buf []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(buf)
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}
