package initialize

import (
	"fc-admin-server-go/global"
	"fc-admin-server-go/models/database"
	"fc-admin-server-go/pkg/config"
	redisutil "fc-admin-server-go/pkg/redis"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func RedisSetUp() *redis.Pool {
	//redis链接池链接	https://phpmianshi.com/?id=4778
	return &redis.Pool{
		MaxIdle:     config.Data.Redis.MaxIdle,     //pool中最大Idle连接数量
		MaxActive:   config.Data.Redis.MaxActive,   //pool中最大分配的连接数量，设为0无限制
		IdleTimeout: config.Data.Redis.IdleTimeout, //idle的时间，超过idle时间连接关闭。设为0 idle的连接不close
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(config.Data.Redis.RedisNetwork, config.Data.Redis.RedisHost)
			if err != nil {
				global.FC_LOGGER.Fatal(fmt.Sprintf("【Redis】链接失败：%v", err)) //redis链接失败,链接失败后可以关闭程序了，所以使用logging.Fatal方法
				return nil, err
			}
			if config.Data.Redis.RedisPassword != "" {
				if _, err := c.Do("AUTH", config.Data.Redis.RedisPassword); err != nil {
					c.Close()
					global.FC_LOGGER.Fatal(fmt.Sprintf("【Redis】链接失败：%v", err)) //redis链接失败,链接失败后可以关闭程序了，所以使用logging.Fatal方法
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				global.FC_LOGGER.Fatal(fmt.Sprintf("【Redis】链接失败：%v", err)) //redis链接失败,链接失败后可以关闭程序了，所以使用logging.Fatal方法
			}
			return err
		},
	}
}

// InitRedisData 启动成功后，同步下role信息
func InitRedisData() {
	roles, err := database.QueryAllRole()
	if err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【Redis】同步数据库Role信息链接Redis失败：%v", err))
	}
	for _, role := range roles {
		err = redisutil.AddRoleInfo(role.ID, role.UpdatedAt)
		if err != nil {
			global.FC_LOGGER.Fatal(fmt.Sprintf("【Redis】同步数据库Role信息失败：%v", err))
			return
		}
	}
}
