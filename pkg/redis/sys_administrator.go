package redisutil

import (
	"encoding/json"
	"fc-admin-server-go/global"
	"fc-admin-server-go/pkg/config"
	"fc-admin-server-go/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"time"
)

var adminKeyName = "AdminInfo:"

// RedisAdminInfo redis中存储的管理员信息
type RedisAdminInfo struct {
	Token         string    `json:"token"`         //管理员token，做唯一登录判断
	AdminStatus   int       `json:"adminStatus"`   //管理员状态，=1正常使用，=2停用 ...
	AdminUpdateAt time.Time `json:"adminUpdateAt"` //管理员信息更新时间，与token中的判断是否需要重新获取管理员信息
	PwdUpdateAt   time.Time `json:"pwdUpdateAt"`   //管理员角色信息密码更新，与token中的判断是否需要重新登录
	RoleIDs       []uint    `json:"roleIDs"`       //管理员角色ID列表
}

// AddAdminInfo redis中添加管理员信息,设置token过期为RefreshExpireTime所设置的时间
func AddAdminInfo(aID uint, appType string, value RedisAdminInfo) error {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	// 将值转换为字节数组
	rAIBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// 设置哈希字段和值
	_, err = conn.Do("HSET", adminKeyName+appType, aID, rAIBytes)
	if err != nil {
		return err
	}
	// 设置过期时间
	_, err = conn.Do("EXPIRE", adminKeyName+appType, int((config.Data.Token.RefreshExpireTime * time.Hour).Seconds()))
	return err
}

// DeleteAdminInfo 删除redis中管理员信息
func DeleteAdminInfo(aID uint, appType string) (err error) {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	_, err = conn.Do("HDEL", adminKeyName+appType, aID)
	return err
}

// GetAdminInfo 获取redis中的指定管理员信息
func GetAdminInfo(aID uint, appType string) (adminInfo RedisAdminInfo, err error) {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	rAIBytes, err := redis.Bytes(conn.Do("HGET", adminKeyName+appType, aID))
	if err != nil {
		return
	}
	err = json.Unmarshal(rAIBytes, &adminInfo) // 将JSON字节数组解析为切片
	return
}

// JwtSetAdminInfoUpdateAt 根据jwt中信息更新redis中管理员信息修改时间(adminUpdateAt管理员信息更新时间)
func JwtSetAdminInfoUpdateAt(adminUpdateAt time.Time, c *gin.Context) (newToken string, err error) {
	tokenString := c.GetHeader("Authorization") //从请求的header中获取token字符串
	rToken, err := util.VerifyAdminToken(tokenString)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		return
	}
	//取token中保存的管理员信息
	claims, err := util.ParseAdminToken(rToken)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		return
	}
	//从redis中获取存储的管理员信息
	appTypeString := c.GetHeader("Type") //从请求的header中获取应用类型字符串，比如Web,H5,iOS,Android
	redisAdminInfo, err := GetAdminInfo(claims.AdminID, appTypeString)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}
		return
	}
	//设置管理员信息更新时间
	redisAdminInfo.AdminUpdateAt = adminUpdateAt
	claims.AdminUpdateAt = adminUpdateAt
	//生成新的token
	newToken, err = util.GenAdminAToken(claims.JwtData, claims.ExpiresAt.Time)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
		return
	}
	//是否做同端唯一登录,是的话更改redis中存储的管理员token
	if config.Data.Token.Unique {
		redisAdminInfo.Token = newToken
	}
	//redis重新存数据
	err = AddAdminInfo(claims.AdminID, appTypeString, redisAdminInfo)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
	}
	return
}

// JwtSetAdminPwdUpdateAt 根据jwt中信息更新redis中管理员密码修改时间(pwdUpdateAt管理员密码更新时间)
func JwtSetAdminPwdUpdateAt(aID uint, appTypeString string, pwdUpdateAt time.Time) (err error) {
	//从redis中获取存储的管理员信息
	redisAdminInfo, err := GetAdminInfo(aID, appTypeString)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
		}
		// redis中未查询到，可能是管理员未登录，所以不用返回错误
		return nil
	}
	//设置管理员密码更新时间
	redisAdminInfo.PwdUpdateAt = pwdUpdateAt
	//redis重新存数据
	err = AddAdminInfo(aID, appTypeString, redisAdminInfo)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
	}
	return
}

// SetAdminInfoUpdateAtStatus 更新redis中管理员信息修改时间(adminUpdateAt管理员信息更新时间) 和 管理员状态
func SetAdminInfoUpdateAtStatus(aID uint, appTypeString string, adminUpdateAt time.Time, status int) (err error) {
	//从redis中获取存储的管理员信息
	redisAdminInfo, err := GetAdminInfo(aID, appTypeString)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}
		// redis中未查询到，可能是管理员未登录，所以不用返回错误
		return nil
	}
	//设置管理员密码更新时间
	redisAdminInfo.AdminUpdateAt = adminUpdateAt
	//管理员状态
	redisAdminInfo.AdminStatus = status
	//redis重新存数据
	err = AddAdminInfo(aID, appTypeString, redisAdminInfo)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
	}
	return
}

// SetAdminInfoRoleIDs 更新redis中管理员角色ID数组
func SetAdminInfoRoleIDs(aID uint, appTypeString string, roleIDs []uint) (err error) {
	//从redis中获取存储的管理员信息
	redisAdminInfo, err := GetAdminInfo(aID, appTypeString)
	if err != nil {
		if err.Error() != "redigo: nil returned" {
			global.FC_LOGGER.Error(fmt.Sprint(err))
			return
		}
		// redis中未查询到，可能是管理员未登录，所以不用返回错误
		return nil
	}
	//设置管理员角色ID数组
	redisAdminInfo.RoleIDs = roleIDs
	//redis重新存数据
	err = AddAdminInfo(aID, appTypeString, redisAdminInfo)
	if err != nil {
		global.FC_LOGGER.Error(fmt.Sprint(err))
	}
	return
}
