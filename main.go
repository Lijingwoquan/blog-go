package main

import (
	"blog/dao/mysql"
	"blog/logger"
	"blog/pkg/snowflake"
	"blog/routers"
	"blog/setting"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func main() {
	//1.加载配置文件
	if err := setting.Init(); err != nil {
		fmt.Println("init setting failed!")
	}
	//2.初始化日志
	if err := logger.Init(viper.GetString("app.mode")); err != nil {
		return
	}
	defer func() {
		if err := zap.L().Sync(); err != nil {
			fmt.Printf(" zap.L().Sync() failed,err:%v", err)
		}
	}()

	//3.初始化mysql
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed! err:%v", err)
		return
	}

	//4.初始化redis
	//if err := redis.Init(); err != nil {
	//	fmt.Printf("init redis failed err:%v", err)
	//	return
	//}
	//defer redis.Close()

	//初始化雪花算法
	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		fmt.Printf("snowflake init failed,err:%v", err)
		return
	}
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 处理 panic，可以记录日志或采取其他措施
				zap.L().Error("Recover from panic in ticker goroutine", zap.Any("panic", r))
			}
		}()

		for range ticker.C {
			// 清理过期的 token
			err := mysql.CleanupInvalidTokens()
			if err != nil {
				zap.L().Error("cleanupExpiredTokens(db) failed", zap.Error(err))
			}
		}
	}()

	//5.注册路由
	r := routers.SetupRouter(viper.GetString("app.mode"))
	err := r.Run(":8080")
	if err != nil {
		fmt.Printf("run server failed,err:%v", err)
		return
	}

}
