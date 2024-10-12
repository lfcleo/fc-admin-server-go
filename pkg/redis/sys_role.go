package redisutil

import (
	"fc-admin-server-go/global"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var roleKeyName = "RoleInfo"

// AddRoleInfos redis中添加多个角色更新时间(存储的是时间戳格式)
func AddRoleInfos(roleIDs []uint, updateTime time.Time) error {
	if len(roleIDs) == 0 {
		return nil
	}
	// 设置redis中角色更新时间
	for _, id := range roleIDs {
		if err := AddRoleInfo(id, updateTime); err != nil {
			return err
		}
	}
	return nil
}

// AddRoleInfo redis中添加角色更新时间(存储的是时间戳格式)
func AddRoleInfo(rID uint, updateTime time.Time) error {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	var unixTime int64
	if updateTime.IsZero() {
		unixTime = 0
	} else {
		unixTime = updateTime.Unix()
	}
	// 设置哈希字段和值
	_, err := conn.Do("HSET", roleKeyName, rID, unixTime)
	return err
}

// GetRoleInfo 获取redis中的指定角色更新时间
func GetRoleInfo(rID uint) (updateTime int64, err error) {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	unixTime, err := redis.Int(conn.Do("HGET", roleKeyName, rID))
	if err != nil {
		return
	}
	updateTime = int64(unixTime)
	return
}

// GetRoleInfos 获取redis中的指定的多个角色信息
func GetRoleInfos(rIDs []uint) (timeUnix []int64, err error) {
	conn := global.FC_REDIS.Get()
	defer conn.Close()

	// 构建参数列表,把roleKeyName添加到第一个元素中
	args := make([]interface{}, len(rIDs)+1)
	args[0] = roleKeyName
	for i, id := range rIDs {
		args[i+1] = id
	}

	values, err := redis.Values(conn.Do("HMGET", args...))
	if err != nil {
		return
	}
	for _, v := range values {
		// 检查每个返回值是否为 nil 或 []byte
		if v == nil {
			continue
		}
		// 尝试将值转换为 []byte
		b, ok := v.([]byte)
		if !ok {
			return nil, fmt.Errorf("redis中解析角色更新时间数据失败1")
		}
		unixNum, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return nil, err
		}
		timeUnix = append(timeUnix, unixNum)
	}
	return
}

// GetAllRoleInfos 获取redis中所有的角色信息
func GetAllRoleInfos() (timeUnix []int64, err error) {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	roleInfos, err := redis.StringMap(conn.Do("HGETALL", roleKeyName))
	if err != nil {
		return
	}
	for _, v := range roleInfos {
		// 检查每个返回值是否为 0
		if v == "" {
			continue
		}
		unixNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		timeUnix = append(timeUnix, unixNum)
	}
	return
}
