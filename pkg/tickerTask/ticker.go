package ticker

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func Init() {
	errCh := make(chan error)
	done := make(chan bool)
	//cleanupInvalidTokensTask
	go func() {
		if err := cleanupInvalidTokensTask(); err != nil {
			errCh <- err
		}
	}()

	//saveVisitedTimesTask
	go func() {
		if err := saveVisitedTimesTask(); err != nil {
			errCh <- err
		}
		done <- true
	}()

	//错误处理
	go func() {
		for err := range errCh {
			zap.L().Error("ticker in pkg happen err:%v", zap.Error(err))
		}
		done <- true
	}()
	go func() {
		<-done
		<-done
		close(errCh)
	}()
}

func cleanupInvalidTokensTask() error {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	for range ticker.C {
		// 清理过期的 token
		err := mysql.CleanupInvalidTokens()
		if err != nil {
			return fmt.Errorf("mysql.CleanupInvalidTokens() failed,err:%v", err)
		}
	}
	return nil
}

func saveVisitedTimesTask() error {
	ticker := time.NewTicker(time.Hour * 1)
	defer ticker.Stop()
	for range ticker.C {
		// 清理过期的 token
		visitedTimesChangedMap, err := redis.GetChangedVisitedTimes()
		if err != nil {
			return fmt.Errorf("redis.GetChangedVisitedTimes() failed,err:%v", err)
		}
		if mysql.SaveVisitedTimes(visitedTimesChangedMap) != nil {
			return fmt.Errorf("mysql.SaveVisitedTimes(visitedTimesChangedMap) failed,err:%v", err)
		}
	}
	return nil
}
