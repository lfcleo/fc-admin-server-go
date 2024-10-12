package middleware

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/response"
	redisUtil "fc-admin-server-go/pkg/redis"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
)

// VerCache redis角色信息鉴权
func VerCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") //从请求的header中获取token字符串
		appTypeString := c.GetHeader("Type")        //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
		//获取token信息
		token, _ := util.VerifyAdminToken(tokenString)
		claims, _ := util.ParseAdminToken(token)

		//取出redis中管理员信息，做判断
		redisAdminInfo, _ := redisUtil.GetAdminInfo(claims.AdminID, appTypeString)
		//如果角色包含超级管理员
		if isOk := util.UintContains(claims.RoleIDs, 1); isOk {
			// 判断超级管理员的角色权限更新时间与redis中角色更新时间比较
			if isOk = judgeAllRedisRoleInfo(redisAdminInfo, claims, c); isOk == false {
				c.Abort()
				return
			}
		} else {
			// 如果角色不包含超级管理员
			// 从redis获取角色信息，判断是否继续下一步
			if isOk = judgeRedisRoleInfo(redisAdminInfo, claims, c); isOk == false {
				c.Abort()
				return
			}
		}
		//设置管理员的ID，供后续方法使用
		c.Set("AID", claims.AdminID)
		c.Next()
	}
}

// 从redis获取指定角色信息，判断
func judgeRedisRoleInfo(redisAdminInfo redisUtil.RedisAdminInfo, claims *util.JwtClaims, c *gin.Context) bool {
	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	//根据管理员的roleIDs，取出redis中保存的角色信息
	redisTimes, err := redisUtil.GetRoleInfos(claims.RoleIDs)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}
		response.Json(response.AuthError, "token参数错误1", c)
		return false
	}
	if len(redisTimes) > 0 {
		//获取角色信息中，最新的更新时间
		redisRoleUpdateAt := getRoleNewUpdateAt(redisTimes)
		//如果有 0 值，代表有角色信息已经删除了，需要管理员重新登录
		if redisRoleUpdateAt == 0 {
			response.Json(response.AuthRoleUpdate, "账户有失效的角色信息，请重新登录！", c)
			return false
		}
		//比较用户token中角色更新时间，和redis保存的role最新更新时间是否一致，不一致退出登录
		if claims.RoleUpdateUnix != redisRoleUpdateAt {
			claims.RoleUpdateUnix = redisRoleUpdateAt
			newToken, err := getTempToken(appTypeString, claims, redisAdminInfo)
			if err != nil {
				response.Json(response.AuthError, "获取管理员信息错误1", c)
				return false
			}
			c.Header("Token", newToken)
			response.Json(response.AuthRoleUpdate, "2-账户角色信息有更新，请重新登录！", c)
			return false
		}
	}
	return true
}

// 从redis获取所有角色信息，判断
func judgeAllRedisRoleInfo(redisAdminInfo redisUtil.RedisAdminInfo, claims *util.JwtClaims, c *gin.Context) bool {
	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	//根据管理员的roleIDs，取出redis中保存的角色信息
	//redisTimes, err := redisUtil.GetAllRoleInfos()
	redisTimes, err := redisUtil.GetRoleInfo(1)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}
		response.Json(response.AuthError, "token参数错误1", c)
		return false
	}
	//if len(redisTimes) > 0 {
	//	//获取角色信息中，最新的更新时间
	//	redisRoleUpdateAt := getRoleNewUpdateAt(redisTimes)
	//	//如果有 0 值，代表有角色信息已经删除了，需要管理员重新登录
	//	if redisRoleUpdateAt == 0 {
	//		response.Json(response.AuthRoleUpdate, "账户有失效的角色信息，请重新登录！", c)
	//		return false
	//	}
	//	log.Println(claims.RoleUpdateUnix, redisRoleUpdateAt)
	//	//比较用户token中角色更新时间，和redis保存的role最新更新时间是否一致，不一致退出登录
	//	if claims.RoleUpdateUnix != redisRoleUpdateAt {
	//		response.Json(response.AuthRoleUpdate, "2-账户角色信息有更新，请重新登录！", c)
	//		return false
	//	}
	//}
	//比较用户token中角色更新时间，和redis保存的role最新更新时间是否一致，不一致退出登录
	if claims.RoleUpdateUnix != redisTimes {
		claims.RoleUpdateUnix = redisTimes
		newToken, err := getTempToken(appTypeString, claims, redisAdminInfo)
		if err != nil {
			response.Json(response.AuthError, "获取管理员信息错误1", c)
			return false
		}
		c.Header("Token", newToken)
		response.Json(response.AuthRoleUpdate, "2-账户角色信息有更新，请重新登录！", c)
		return false
	}
	return true
}
