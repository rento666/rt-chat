package thirdapi

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"regexp"
	"rt-chat/internal/result"
	"rt-chat/pkg/tools/ttrans"
)

// Api60s 每天60秒看世界api
// @Tags 每天60秒看世界
// Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {string} json "{"id": "str"}"
// @Router /v1/tool/60s [get]
func Api60s(c *gin.Context) {
	r := get60s()
	data := ttrans.Transformation(r)
	code := http.StatusNotFound
	msg := "获取原始数据成功"
	var res interface{}
	res = data
	// 找到data[0],此时res为data[0],类型为json
	if v1, ok1 := data["data"].([]interface{}); ok1 {
		code = http.StatusOK
		msg = "获取data数据成功"
		res = v1[0]
	}
	// 找到content,此时res类型为string
	if v2, ok2 := res.(map[string]interface{}); ok2 {
		code = http.StatusOK
		msg = "获取content数据成功"
		res = v2["content"]
	}
	if v3, ok3 := res.(string); ok3 {
		var arr []string
		// 将content数据进行正则匹配
		compile := regexp.MustCompile(`<p\s+data-pid=[^<>]+>([^<>]+)</p>`)
		// 将标签匹配<></>
		sr := regexp.MustCompile(`<[^<>]+>`)
		// 匹配 &#34; => "
		ascii := regexp.MustCompile(`&#(\d+);`)
		// 拿到p标签的汉字
		f := compile.FindAllString(v3, -1)
		for _, v := range f {
			temp := sr.ReplaceAllString(v, "")
			arr = append(arr, ascii.ReplaceAllString(temp, ""))
		}
		res = arr
	}
	c.JSON(http.StatusOK, result.Json(code, msg, res))
}

// Get60s 获取 “每天 60 秒读懂世界” 数据
/*
 * @return string型数据
 */
func get60s() []byte {
	url := "https://www.zhihu.com/api/v4/columns/c_1261258401923026944/items?limit=1"
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	res, _ := ttrans.ZhToUnicode(content)
	return res
}
