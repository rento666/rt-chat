package tcaptcha

import (
	"encoding/json"
	"fmt"
	"github.com/wenlng/go-captcha/captcha"
	"math/rand"
	"net/http"
	"rt-chat/internal/global"
	"rt-chat/pkg/tools/tredis"
	"rt-chat/pkg/tools/ttrans"
	"strconv"
	"strings"
	"time"
)

type Ct struct {
	ImgBase64   string
	ThumbBase64 string
	CaptchaKey  string
}
type Che struct {
	dots string
	key  string
}

func CheckCode(dots string, key string) (string, bool, int) {
	// 先声明返回体，默认为失败(false、412 状态码)
	msg := "人机验证失败"
	isTrue := false
	code := http.StatusPreconditionFailed
	// 传递的值不为空则执行if内的语句
	if dots != "nil" || key != "nil" {
		if isTrue = CheckCaptcha(dots, key); isTrue {
			code = http.StatusOK
			msg = "人机验证成功"
		}
	}
	return msg, isTrue, code
}

/**
生成验证码
	需要根据key去redis中获取存储的dots
为了防止同一ip恶意申请生成验证码，使用时需添加判断，例如：如果在半小时内连续申请1k次，则进行短暂拉黑（1天）处理
*/

// CreateCode 创建二维码，返回二维码主图、缩略图、key(redis)
func CreateCode() *Ct {
	// 单例模式得到验证码
	capt := captcha.GetCaptcha()
	// 生成验证码
	dots, b64, tb64, key, err := capt.Generate()
	if err != nil {
		// "message": "GenCaptcha err",
		return nil
	}
	// 将文本位置验证数据（等价于密钥）存入redis
	bt, _ := json.Marshal(dots)
	value := ttrans.Bytes2String(bt)
	tredis.SetString(global.REDIS, key, value, time.Minute*5)
	// 主图base64、缩略图base64、唯一key
	ct := new(Ct)
	ct.ImgBase64 = b64
	ct.ThumbBase64 = tb64
	ct.CaptchaKey = key
	return ct
}

// CheckCaptcha 检查二维码
/**
@param dots string
@param key string
@return bool true为通过验证
*/
func CheckCaptcha(dots string, key string) bool {

	if dots == "" || key == "" {
		// "message": "dots or key param is empty",
		return false
	}

	cacheData := tredis.GetString(global.REDIS, key)
	if cacheData == "" {
		// "message": "illegal key",
		return false
	}
	src := strings.Split(dots, ",")

	var dct map[int]captcha.CharDot
	if err := json.Unmarshal([]byte(cacheData), &dct); err != nil {
		// "message": "illegal key",
		return false
	}

	chkRet := false
	if (len(dct) * 2) == len(src) {
		for i, dot := range dct {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[k]), 64)

			// 检测点位置
			// chkRet = captcha.CheckPointDist(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))

			// 校验点的位置,在原有的区域上添加额外边距进行扩张计算区域,不推荐设置过大的padding
			// 例如：文本的宽和高为30，校验范围x为10-40，y为15-45，此时扩充5像素后校验范围宽和高为40，则校验范围x为5-45，位置y为10-50
			chkRet = captcha.CheckPointDistWithPadding(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), 5)
			if !chkRet {
				break
			}
		}
	}

	if chkRet {
		// 通过校验，把redis中的键值对删除
		_, err := tredis.DelByKey(global.REDIS, key)
		if err != nil {
			// 未正常删除键值对
		}
		return true
	}

	return false
}

// GenerateVerificationCode 生成随机数字验证码
func GenerateVerificationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano())) // 设置随机种子

	verificationCode := ""
	for i := 0; i < 6; i++ {
		verificationCode += strconv.Itoa(rand.Intn(10)) // 生成随机数字，范围是0-9的数字
	}

	return verificationCode
}
