package middleware

import (
	"errors"
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/request"
	"fc-admin-server-go/models/response"
	casbinUtil "fc-admin-server-go/pkg/casbin"
	"fc-admin-server-go/pkg/config"
	redisUtil "fc-admin-server-go/pkg/redis"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"reflect"
	"strconv"
	"strings"
)

var defaultPath = "/v1/token/"

// DefaultPermissionAPIs 路由白名单，访问此路由不做权限判断
func DefaultPermissionAPIs() []*request.DefaultAPIModel {
	return []*request.DefaultAPIModel{
		{Path: "auth/refresh", Method: "POST"},
		{Path: "auth/menu", Method: "POST"},
		{Path: "auth/logout", Method: "POST"},
		{Path: "auth/update/info", Method: "POST"},
		{Path: "auth/update/pwd", Method: "POST"},
		{Path: "auth/info", Method: "POST"},
		{Path: "auth/dict", Method: "POST"},
		{Path: "auth/upload", Method: "POST"},
	}
}

// ExistsInArray 检查是否是默认路由
func ExistsInArray(da *request.DefaultAPIModel) bool {
	for _, v := range DefaultPermissionAPIs() {
		if reflect.DeepEqual(v, da) {
			return true
		}
	}
	return false
}

// VerToken 管理员token鉴权
func VerToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") //从请求的header中获取token字符串
		if tokenString == "" {                      //如果没有的话，直接返回token错误
			response.Json(response.AuthError, "no token", c)
			c.Abort()
			return
		} else {
			//获取token信息
			token, err := util.VerifyAdminToken(tokenString)
			if err != nil {
				//判断token是否过期,过期返回402，需重新获取token（如果是刷新token，额外处理）
				if errors.Is(err, jwt.ErrTokenExpired) {
					//判断，如果是刷新token路由,放行。方法中有判断
					if c.Request.URL.String() == "/v1/token/auth/refresh" {
						c.Next()
						return
					}
					response.Json(response.AuthExpire, "token expire", c)
					c.Abort()
					return
				}
				global.FC_LOGGER.Error(fmt.Sprint(err))
				response.Json(response.AuthError, "verify token error 2", c)
				c.Abort()
				return
			}
			//如果token没有通过校验，返回401
			if !token.Valid {
				response.Json(response.AuthError, "token Valid false", c)
				c.Abort()
				return
			}
			claims, err := util.ParseAdminToken(token)
			if err != nil { //token校验失败，返回错误信息
				global.FC_LOGGER.Error(fmt.Sprint(err))
				response.Json(response.AuthError, "token error", c)
				c.Abort()
				return
			} else {
				// 从redis获取管理员信息，判断是否继续下一步
				_, isOk := judgeRedisAdminInfo(claims, tokenString, c)
				if isOk == false {
					c.Abort()
					return
				}
				//如果角色不包含超级管理员
				if isOk := util.UintContains(claims.RoleIDs, 1); isOk == false {
					// casbin判断用户权限
					if isOk = judgeCasbinRoleInfo(claims, c); isOk == false {
						c.Abort()
						return
					}
				}
				//设置管理员的ID，供后续方法使用
				c.Set("AID", claims.AdminID)
				c.Next()
			}
		}
	}
}

// 从redis获取管理员信息，判断token,管理员状态，管理员密码更新时间，管理员角色信息，判断管理员信息的更新时间
func judgeRedisAdminInfo(claims *util.JwtClaims, tokenString string, c *gin.Context) (redisAdminInfo redisUtil.RedisAdminInfo, isOk bool) {
	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	//取出redis中报错的管理员信息，做判断
	redisAdminInfo, err := redisUtil.GetAdminInfo(claims.AdminID, appTypeString)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}
		response.Json(response.AuthError, "管理员信息错误！", c)
		return redisAdminInfo, false
	}
	//是否做同端唯一登录,比较token
	if config.Data.Token.Unique {
		if redisAdminInfo.Token != tokenString {
			response.Json(response.AuthError, "令牌错误", c)
			return redisAdminInfo, false
		}
	}
	//判断管理员状态是否可用
	if redisAdminInfo.AdminStatus != 1 {
		response.Json(response.AuthStatusError, "当前管理员状态不可用，请联系超级管理员！", c)
		return redisAdminInfo, false
	}
	//判断管理员密码的更新时间，是否需要登录
	if claims.PwdUpdateAt.Equal(redisAdminInfo.PwdUpdateAt) == false {
		response.Json(response.AuthPwdUpdate, "密码已修改，请重新登录！", c)
		return redisAdminInfo, false
	}
	//判断管理员角色信息，是否需要登录
	if util.UintArraysEqual(claims.RoleIDs, redisAdminInfo.RoleIDs) == false {
		claims.RoleIDs = redisAdminInfo.RoleIDs
		newToken, err := getTempToken(appTypeString, claims, redisAdminInfo)
		if err != nil {
			response.Json(response.AuthError, "获取管理员信息错误1", c)
			return redisAdminInfo, false
		}
		c.Header("Token", newToken)
		response.Json(response.AuthRoleUpdate, "1-账户角色信息有更新，请重新登录！", c)
		return redisAdminInfo, false
	}
	//判断管理员信息的更新时间，是否需要重新获取用户信息（header返回新token）
	if claims.AdminUpdateAt.Equal(redisAdminInfo.AdminUpdateAt) == false {
		claims.AdminUpdateAt = redisAdminInfo.AdminUpdateAt
		newToken, err := getTempToken(appTypeString, claims, redisAdminInfo)
		if err != nil {
			response.Json(response.AuthError, "获取管理员信息错误1", c)
			return redisAdminInfo, false
		}
		c.Header("Token", newToken)
		response.Json(response.AuthInfoUpdate, "please update admin info", c)
		return redisAdminInfo, false
	}
	return redisAdminInfo, true
}

// 获取新token
func getTempToken(appTypeString string, claims *util.JwtClaims, redisAdminInfo redisUtil.RedisAdminInfo) (newToken string, err error) {
	//生成新的token
	newToken, err = util.GenAdminAToken(claims.JwtData, claims.ExpiresAt.Time)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		return
	}
	//是否做同端唯一登录,是的话更改redis中存储的管理员token
	if config.Data.Token.Unique {
		redisAdminInfo.Token = newToken
		//redis重新存数据
		err = redisUtil.AddAdminInfo(claims.AdminID, appTypeString, redisAdminInfo)
		if err != nil {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}
	}
	return
}

// 获取角色中最新的更新时间
func getRoleNewUpdateAt(timeUnix []int64) int64 {
	newTime := timeUnix[0]
	//如果有 0 值，代表此角色信息已经删除了，需要管理员重新登录
	if newTime == 0 {
		return 0
	}
	for _, t := range timeUnix {
		//如果有 0 值，代表此角色信息已经删除了，需要管理员重新登录
		if t == 0 {
			return 0
		}
		if t > newTime {
			newTime = t
		}
	}
	return newTime
}

// 从casbin获取角色信息，判断
func judgeCasbinRoleInfo(claims *util.JwtClaims, c *gin.Context) bool {
	//获取请求的PATH
	obj := strings.TrimPrefix(c.Request.URL.Path, defaultPath)
	// 获取请求方法
	act := c.Request.Method
	//如果是否在路由白名单中
	dAPi := request.DefaultAPIModel{
		Path:   obj,
		Method: act,
	}
	if isOk := ExistsInArray(&dAPi); isOk {
		return true
	}
	// casbin 检测路由角色权限
	for _, rID := range claims.RoleIDs {
		sub := strconv.Itoa(int(rID))
		success, err := casbinUtil.CasbinServiceApp.CanAccess(sub, obj, act)
		if err != nil {
			response.Json(response.AuthPoor, "无权限访问或操作!", c)
			return false
		}
		if success {
			return true
		}
	}
	response.Json(response.AuthPoor, "无权限访问或操作!", c)
	return false
}
