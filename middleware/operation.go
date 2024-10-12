package middleware

import (
	"bytes"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"reflect"
	"strings"
	"sync"
)

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

// DefaultIgnoreAPIs 忽略路由名单，访问此路由不添加操作日志中
func defaultIgnoreAPIs() []*request.DefaultAPIModel {
	return []*request.DefaultAPIModel{
		{Path: "/v1/token/auth/refresh", Method: "POST"},
		{Path: "/v1/token/auth/menu", Method: "POST"},
		{Path: "/v1/token/auth/info", Method: "POST"},
		{Path: "/v1/token/auth/dict", Method: "POST"},
		{Path: "/v1/token/admin/list", Method: "POST"},
		{Path: "/v1/token/role/list", Method: "POST"},
		{Path: "/v1/token/role/apis", Method: "POST"},
		{Path: "/v1/token/role/menus", Method: "POST"},
		{Path: "/v1/token/menu/list", Method: "POST"},
		{Path: "/v1/token/api/list", Method: "POST"},
		{Path: "/v1/token/api/list/all", Method: "POST"},
		{Path: "/v1/token/dict/type/list", Method: "POST"},
		{Path: "/v1/token/dict/data/list", Method: "POST"},
		{Path: "/v1/token/log/list", Method: "POST"},
		{Path: "/v1/token/log/delete", Method: "POST"},
	}
}

// ExistsInArray 检查是否是忽略路由
func existsInArray(da *request.DefaultAPIModel) bool {
	for _, v := range defaultIgnoreAPIs() {
		if reflect.DeepEqual(v, da) {
			return true
		}
	}
	return false
}

// OperationLog 添加操作记录
func OperationLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//如果是否在路由白名单中
		dAPi := request.DefaultAPIModel{
			Path:   c.Request.URL.Path, //请求的PATH
			Method: c.Request.Method,   //请求方法
		}
		if isOk := existsInArray(&dAPi); isOk {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()

		var adminID uint
		// 判断管理员ID，在没有token中间件的路由，根据手机号查找ID
		if strings.Contains(c.Request.URL.Path, "token") {
			adminID = c.MustGet("AID").(uint)
		} else {
			var aAuth request.Auth
			if err := c.ShouldBindBodyWithJSON(&aAuth); err != nil {
				//global.FC_LOGGER.Error(fmt.Sprint(err))
				return
			}
			//正则验证手机号/邮箱是否符合要求
			isMobile := util.ValidateMobileNumber(aAuth.Username)
			isEmail := util.ValidateEmail(aAuth.Username)
			if isMobile == false && isEmail == false {
				//global.FC_LOGGER.Error(aAuth.Username + "：手机号/邮箱格式错误")
				return
			}
			//根据手机号/邮箱查询管理员信息
			queryKey := "mobile"
			if isEmail {
				queryKey = "email"
			}
			//数据库查询管理员信息
			administrator, err := database.FindAdministratorByKey(queryKey, aAuth.Username, "RolesData")
			if err != nil {
				global.FC_LOGGER.Error(fmt.Sprint(err))
				return
			}
			if administrator == nil {
				//global.FC_LOGGER.Error(aAuth.Username + "：手机号/邮箱数据库不存在")
				return
			}
			adminID = administrator.ID
		}
		oLog := database.SysOperationLog{
			IP:              c.ClientIP(),
			Method:          c.Request.Method,
			Path:            c.Request.URL.Path,
			Agent:           c.Request.UserAgent(),
			Response:        writer.body.String(),
			Code:            c.Writer.Status(),
			AdministratorID: adminID,
		}

		//判断，如果是上传文件接口，上传的是form-data
		if c.ContentType() == "multipart/form-data" {
			oLog.Request = "{\"文件上传\":\"[文件]\"}"
		} else {
			oLog.Request = string(body)
		}
		if err = oLog.CreateOperationLog(); err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}

	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
