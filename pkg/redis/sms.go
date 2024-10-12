package redisutil

import (
	"fc-admin-server-go/global"
	"github.com/gomodule/redigo/redis"
	"time"
)

// AddSmsCode redis中手机号与对应的验证码，有效时间5分钟
func AddSmsCode(mobile, code string) error {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	//添加redis中，这个是没有添加过期时间的	uint转string
	_, err := conn.Do("SET", mobile, code, "EX", int((5 * time.Minute).Seconds()))
	return err
}

// GetSmsCode redis根据手机号查找验证码
func GetSmsCode(mobile string) (string, error) {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	tokenStr, err := redis.String(conn.Do("GET", mobile))
	return tokenStr, err
}

// DelSmsCode redis根据手机号删除保存的验证码信息
func DelSmsCode(mobile string) error {
	conn := global.FC_REDIS.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", mobile)
	return err
}
