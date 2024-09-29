package ticker

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	cleanInvalidToken       = time.Hour * 24
	saveVisitedTimes        = time.Hour * 24
	cleanLowFrequentKeyword = time.Hour * 24 * 30
	taskCount               = 2
)

var (
	wg    sync.WaitGroup
	errCh = make(chan error, taskCount)
)

func Init() {
	wg.Add(taskCount + 1) // 为所有任务和错误处理 goroutine 加 1

	taskList := []func() error{saveVisitedTimesTask, cleanLowFrequentKeywordTask}

	// 启动错误处理goroutine
	go func() {
		defer wg.Done()
		handleErrors()
	}()

	// 启动任务
	for _, task := range taskList {
		go func(t func() error) {
			defer wg.Done()
			if err := runTask(t); err != nil {
				errCh <- err
			}
		}(task)
	}
	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(errCh) // 所有任务完成后关闭错误通道
	}()
}

func runTask(task func() error) error {
	return task()
}

func handleErrors() {
	for err := range errCh {
		// 使用你的日志库记录错误，这里用zap作为示例
		zap.L().Error("Task happen error", zap.Error(err))
	}
}

func cleanupInvalidTokensTask() error {
	ticker := time.NewTicker(cleanInvalidToken)
	defer ticker.Stop()
	for range ticker.C {
		// 清理过期的 token
		err := mysql.CleanupInvalidTokens()
		if err != nil {
			return err
		}
	}
	return nil
}

func saveVisitedTimesTask() error {
	ticker := time.NewTicker(saveVisitedTimes)
	defer ticker.Stop()
	for range ticker.C {
		// 得到浏览次数
		visitedTimesChangedMap, err := redis.GetAndClearChangedVisitedTimes()
		if err != nil {
			return err
		}
		if mysql.SaveVisitedTimes(visitedTimesChangedMap) != nil {
			return err
		}
	}
	return nil
}

func cleanLowFrequentKeywordTask() error {
	ticker := time.NewTicker(cleanLowFrequentKeyword)
	defer ticker.Stop()
	for range ticker.C {
		err := redis.CleanLowerKeywordsZsetEveryMonth()
		if err != nil {
			return err
		}

	}
	return nil
}
