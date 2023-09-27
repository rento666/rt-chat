package tencrypt

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encrypt MD5加密
func MD5Encrypt(string string) string {
	h := md5.New()
	h.Write([]byte(string)) // 需要加密的字符串为
	return hex.EncodeToString(h.Sum(nil))
}

// CodeSalting 用户密码加盐，需要密码、盐、用户id
func CodeSalting(password, salt string, uid uint) string {
	// 随机截取盐的长度（这里的盐是生成的uuid，长度大于10位）
	num := uid % 10
	return password + salt[:num]
}
