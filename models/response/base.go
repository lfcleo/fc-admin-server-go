package response

import (
	"encoding/base64"
	"encoding/json"
	"fc-admin-server-go/global"
	"fc-admin-server-go/pkg/config"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	SUCCESS         = 200 //成功响应请求
	ERROR           = 500 //错误响应请求
	AuthError       = 401 //token无效,需重新登录。
	AuthExpire      = 402 //token过期,需重新获取。
	AuthStatusError = 403 //账号被冻结,联系超级管理员处理。
	AuthPoor        = 404 //无权限操作API
	AuthInfoUpdate  = 410 //账户信息有更新，需重新获取管理员信息。
	AuthRoleUpdate  = 411 //账户角色信息有更新，需重重新登录。
	AuthPwdUpdate   = 412 //账户角色信息密码更新，需重重新登录。
)

// MsgFlags 编码消息
var MsgFlags = map[int]string{
	SUCCESS:         "success",
	ERROR:           "fail",
	AuthError:       "token 无效",
	AuthExpire:      "token 过期",
	AuthStatusError: "当前账号被冻结,联系超级管理员处理",
	AuthInfoUpdate:  "当前账号信息有更新，需重新获取管理员信息。",
	AuthRoleUpdate:  "当前账号角色信息有更新，需重新菜单和权限信息。",
	AuthPwdUpdate:   "当前账号角色信息密码更新，需重重新登录。",
}

// Response 数据返回信息的model，格式如下
type Response struct {
	Code    int         `json:"code"`    //自定义编码
	Message string      `json:"message"` //自定义消息
	Data    interface{} `json:"data"`    //返回的数据
}

// Json 返回数据
func Json(code int, data interface{}, c *gin.Context) {
	msg := ""
	if code != 200 {
		msg = fmt.Sprint(data)
	}

	//判断如果是请求加密，返回也加密
	encryptionString := c.GetHeader("Encryption")
	if encryptionString == "true" {
		ivKey := c.MustGet("ivKey").(string)
		eData, err := encryptData(data, ivKey)
		if err != nil {
			c.JSON(SUCCESS, &Response{
				Code:    ERROR,
				Message: "服务器加密数据struct转化错误",
				Data:    "服务器加密数据struct转化错误",
			})
			return
		}
		data = base64.StdEncoding.EncodeToString(eData)
		c.Header("Encryption", "true")
	}

	// 返回结果数据
	c.JSON(SUCCESS, &Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// encryptData 对返回数据加密
func encryptData(data interface{}, ivKey string) (encryptData []byte, err error) {
	// 判断data类型，转为[]byte类型
	var bytes []byte
	switch v := data.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	case gin.H:
		bytes, err = json.Marshal(data)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}

	default:
		bytes, err = json.Marshal(data)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}
	}

	encryptData, err = util.AesCBCPkcs7Encrypt(bytes, config.Data.Server.RequestSign, ivKey)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		return
	}
	return
}
