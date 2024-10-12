package global

import (
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	FC_DB     *gorm.DB
	FC_REDIS  *redis.Pool
	FC_LOGGER *zap.Logger
)
