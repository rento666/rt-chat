package ttrans

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// String2Bytes string转[]byte
func String2Bytes(s string) []byte {
	return []byte(s)
}

// Bytes2String []byte转string
func Bytes2String(b []byte) string {
	//return *(*string)(unsafe.Pointer(&b))
	return string(b)
}

// SplitInString 分割字符串
func SplitInString(str string, sep string) []string {
	return strings.Split(str, sep)
}

// String2uInt8 string转uint8
func String2uInt8(s string) uint8 {
	ints, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}
	return uint8(ints)
}

// ZhToUnicode 转字符 : "str" => str
func ZhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// Transformation []byte 转 map(key,value)
func Transformation(body []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(string(body)), &result)
	return result
}

func DurationToMinute(time time.Duration) time.Duration {
	return time / 60000000000
}
