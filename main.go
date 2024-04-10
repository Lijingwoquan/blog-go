package main

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/logger"
	"blog/pkg/snowflake"
	"blog/routers"
	"blog/setting"
	"database/sql"
	"errors"
	"fmt"
	redis2 "github.com/go-redis/redis"
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
	//Wait for MySQL
	if err := waitForMySQL(); err != nil {
		fmt.Printf("Error waiting for MySQL: %v\n", err)
		return
	}

	//Wait for Redis
	if err := waitForRedis(); err != nil {
		fmt.Printf("Error waiting for Redis: %v\n", err)
		return
	}

	//3.初始化mysql
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failed! err:%v", err)
		return
	}

	//4.初始化redis
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failed err:%v", err)
		return
	}
	defer redis.Close()

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

func waitForMySQL() error {
	dsn := "root:liuzihao520@tcp(mysql)/blog"
	for i := 0; i < 30; i++ { // 尝试等待的次数，可以根据实际情况调整
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			fmt.Printf("Error opening MySQL connection: %v\n", err)
		} else {
			defer db.Close()
			if err := db.Ping(); err == nil {
				fmt.Println("MySQL is ready!")
				return nil
			} else {
				fmt.Printf("Error pinging MySQL: %v\n", err)
			}
		}
		fmt.Printf("Waiting for MySQL (attempt %d)...\n", i+1)
		time.Sleep(5 * time.Second)
	}
	return errors.New("timed out waiting for MySQL")
}

func waitForRedis() error {
	client := redis2.NewClient(&redis2.Options{
		Addr:     "redis:6379",
		Password: "", // 如果有密码，请设置密码
		DB:       0,
	})

	for {
		err := client.Ping().Err()
		if err == nil {
			fmt.Println("Redis is ready!")
			return nil
		}

		fmt.Println("Error pinging Redis:", err)
		fmt.Println("Waiting for Redis...")
		time.Sleep(2 * time.Second)
	}
}
