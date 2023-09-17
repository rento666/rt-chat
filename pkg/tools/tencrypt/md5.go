package tencrypt

import (
	"crypto/md5"
	"encoding/hex"
)

/**
  md5加密
  @author Bill
*/

func MD5Encrypt(string string) string {
	h := md5.New()
	h.Write([]byte(string)) // 需要加密的字符串为
	return hex.EncodeToString(h.Sum(nil))
}
