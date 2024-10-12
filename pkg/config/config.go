package config

import (
	"time"
)

var Data = &Config{}

type Config struct {
	Name     string   `yaml:"name"`
	Version  string   `yaml:"version"`
	Debug    bool     `yaml:"debug"`
	Server   Server   `yaml:"server"`
	Token    Token    `yaml:"token"`
	Zap      Zap      `yaml:"zap"`
	Database Database `yaml:"database"`
	Redis    Redis    `yaml:"redis"`
}

type Server struct {
	HttpPort        int           `yaml:"httpPort"`
	DomainName      string        `yaml:"domainName"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	RequestSign     string        `yaml:"requestSign"`
	PasswordSign    string        `yaml:"passwordSign"`
	PageSize        int           `yaml:"pageSize"`
	RuntimeRootPath string        `yaml:"runtimeRootPath"`
	ImageSavePath   string        `yaml:"imageSavePath"`
	ImageMaxSize    int           `yaml:"imageMaxSize"`
	ImageAllowExts  []string      `yaml:"imageAllowExts"`
	VideoSavePath   string        `yaml:"videoSavePath"`
	VideoMaxSize    int           `yaml:"videoMaxSize"`
	VideoAllowExts  []string      `yaml:"videoAllowExts"`
	ApkSavePath     string        `yaml:"apkSavePath"`
	ApkAllowExt     string        `yaml:"apkAllowExt"`
	AppStoreUrl     string        `yaml:"appStoreUrl"`
	TimeFormat      string        `yaml:"timeFormat"`
}

type Token struct {
	Secret                string        `yaml:"secret"`
	AccountExpireTime     time.Duration `yaml:"accountExpireTime"`
	RefreshExpireTime     time.Duration `yaml:"refreshExpireTime"`
	RefreshAutoExpireTime time.Duration `yaml:"refreshAutoExpireTime"`
	Unique                bool          `yaml:"unique"`
}

type Zap struct {
	InfoFilename  string `yaml:"infoFilename"`
	ErrorFilename string `yaml:"errorFilename"`
	PanicFilename string `yaml:"panicFilename"`
	FatalFilename string `yaml:"fatalFilename"`
	MaxSize       int    `yaml:"maxSize"`
	MaxBackups    int    `yaml:"maxBackups"`
	MaxAge        int    `yaml:"maxAge"`
}

type Database struct {
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Host            string `yaml:"host"`
	Name            string `yaml:"name"`
	TablePrefix     string `yaml:"tablePrefix"`
	MaxIdleConnects int    `yaml:"maxIdleConns"`
	MaxOpenConnects int    `yaml:"maxOpenConns"`
}

type Redis struct {
	RedisNetwork  string        `yaml:"redisNetwork"`
	RedisHost     string        `yaml:"redisHost"`
	RedisPassword string        `yaml:"redisPassword"`
	MaxIdle       int           `yaml:"maxIdle"`
	MaxActive     int           `yaml:"maxActive"`
	IdleTimeout   time.Duration `yaml:"idleTimeout"`
}
