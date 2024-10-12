package middleware

import (
	"bytes"
	"encoding/json"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	"fc-admin-server-go/pkg/config"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"time"
)

// appType类型
var validAppType = map[string]bool{
	"Web": true,
	"H5":  true,
}

// ReqEncrypt 加密请求参数结构体
type ReqEncrypt struct {
	request.Base
	EncryptData string `json:"encryptData"` //签名
}

// VerEncrypt 验证数据加密
func VerEncrypt() gin.HandlerFunc {
	return func(c *gin.Context) {

		//判断请求header中的appType是否属于validAppType
		appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
		if !validAppType[appTypeString] {
			response.Json(response.ERROR, "header type error", c)
			c.Abort()
			return
		}

		//判断请求参数是否加密
		encryptionString := c.GetHeader("Encryption")
		if encryptionString == "true" {
			if ok := analysisEncryptionData(c); ok == false {
				c.Abort()
				return
			}
		}

		//通过校验，进行下一步
		c.Next()
	}
}

// analysisEncryptionData 解析加密数据
func analysisEncryptionData(c *gin.Context) bool {

	//判断，如果是上传文件接口，上传的是form-data
	if c.ContentType() == "multipart/form-data" {
		timestamp := c.PostForm("timestamp") //从表单中查询参数
		//时间戳字符串转int64
		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			response.Json(response.ERROR, "解析时间戳错误", c)
			return false
		}
		//判断时间戳是否在1分钟的请求
		if isOk := judgeTimestamp(timestampInt); isOk == false {
			response.Json(response.ERROR, "超时错误", c)
			return false
		}
		encryptData := c.PostForm("encryptData") //从表单中查询参数

		ivKey := strconv.FormatInt(timestampInt*1000, 10)
		panPwdBytes, err := util.AesCBCPkcs7Decrypt(encryptData, config.Data.Server.RequestSign, ivKey)
		if err != nil || len(panPwdBytes) == 0 {
			response.Json(response.ERROR, "加密错误", c)
			return false
		}
		//设置解析后的参数
		c.Request.PostForm.Set("director", string(panPwdBytes))
		//设置ivKey,为后续返回加密使用
		c.Set("ivKey", ivKey)
	} else {

		data, err := c.GetRawData()
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			response.Json(response.ERROR, "加密参数错误", c)
			return false
		}

		var re ReqEncrypt
		// 把请求参数转化为ReqEncrypt结构体,如果转换失败，拒绝此次请求
		if err := json.Unmarshal(data, &re); err != nil {
			response.Json(response.ERROR, "加密校验错误", c)
			return false
		}
		//判断时间戳是否在1分钟的请求
		if isOk := judgeTimestamp(re.Timestamp); isOk == false {
			response.Json(response.ERROR, "超时处理", c)
			return false
		}

		ivKey := strconv.FormatInt(re.Timestamp*1000, 10)
		panPwdBytes, err := util.AesCBCPkcs7Decrypt(re.EncryptData, config.Data.Server.RequestSign, ivKey)
		if err != nil || len(panPwdBytes) == 0 {
			response.Json(response.ERROR, "加密错误", c)
			return false
		}

		//设置ivKey,为后续返回加密使用
		c.Set("ivKey", ivKey)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(panPwdBytes))
	}
	return true
}

// 判断用户请求的时间是否在1分钟内，1分钟内的请求代表有效请求
func judgeTimestamp(timestampInt int64) bool {
	//判断用户请求的时间戳是否在允许的时间之内
	timeST := time.Unix(timestampInt, 0) //用户请求的时间戳转化为时间
	nowTime := time.Now()                //当前时间
	t := nowTime.Sub(timeST)             //计算当前时间与用户请求时间的差值
	tM := t.Minutes()                    //差值转化为时间
	if tM > 1 {                          //如果两个时间差大于1分钟，则代表不能使用
		return false
	}
	return true
}
