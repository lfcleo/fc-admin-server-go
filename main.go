package main

import (
	"context"
	"fc-admin-server-go/global"
	"fc-admin-server-go/initialize"
	"fc-admin-server-go/pkg/config"
	"fc-admin-server-go/routers"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("Hello, Http服务正在启动中...")

	initialize.ConfigSetUp()                    //初始化配置文件
	global.FC_LOGGER = initialize.ZapLogSetUp() //初始化打印日志
	global.FC_DB = initialize.GormSetUp()       //设置数据库
	initialize.InitSQLData()                    //设置数据库初始化数据
	global.FC_REDIS = initialize.RedisSetUp()   //设置redis
	initialize.InitRedisData()                  //设置redis初始化数据
	defer func() {
		// 程序结束前关闭数据库链接
		db, _ := global.FC_DB.DB()
		db.Close()
		// 程序结束前关闭打印日志
		global.FC_LOGGER.Sync()
	}()

	router := routers.InitRouter() //初始化路由
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Data.Server.HttpPort), //设置端口号
		Handler:        router,                                          //http句柄，实质为ServeHTTP，用于处理程序响应HTTP请求
		ReadTimeout:    config.Data.Server.ReadTimeout * time.Second,    //允许读取的最大时间
		WriteTimeout:   config.Data.Server.WriteTimeout * time.Second,   //允许写入的最大时间
		MaxHeaderBytes: 1 << 20,                                         //请求头的最大字节数
	}
	// 使用 http.Server - Shutdown() 优雅的关闭http服务
	go func() {
		if err := s.ListenAndServe(); err != nil {
			global.FC_LOGGER.Info(fmt.Sprintf("【Http服务】监听服务错误: %s", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("结束Http服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		global.FC_LOGGER.Fatal(fmt.Sprintf("【Http服务】服务关闭错误：%s", err))
	}
	log.Println("程序服务关闭退出")
}
