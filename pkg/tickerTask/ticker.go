package ticker

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"go.uber.org/zap"
	"time"
)

func Init() {
	go cleanupInvalidTokensTask()
}

func cleanupInvalidTokensTask() error {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()
	for range ticker.C {
		// 清理过期的 token
		err := mysql.CleanupInvalidTokens()
		if err != nil {
			zap.L().Error("mysql.CleanupInvalidTokens() failed", zap.Error(err))
			return err
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
			zap.L().Error("redis.GetChangedVisitedTimes() failed", zap.Error(err))
			return err
		}
		if mysql.SaveVisitedTimes(visitedTimesChangedMap) != nil {
			zap.L().Error("mysql.SaveVisitedTimes(visitedTimesChangedMap) failed", zap.Error(err))
			return err
		}
	}
	return nil
}
