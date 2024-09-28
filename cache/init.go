package cache

import (
	"go.uber.org/zap"
	"sync"
)

const (
	tickerTaskCount = 3
)

func Init() {
	var wg sync.WaitGroup
	errCh := make(chan error, tickerTaskCount) // 缓冲通道，避免goroutine泄露

	// getDataAboutIndex
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getDataAboutIndex(); err != nil {
			errCh <- err
		}
	}()

	//getEssayList
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getEssayList(); err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := getMaliciousMap(); err != nil {
			errCh <- err
		}
	}()

	// 错误处理
	go func() {
		wg.Wait()
		close(errCh) // 所有任务完成后关闭errCh
		for err := range errCh {
			zap.L().Error("Error in cache Init", zap.Error(err))
		}
	}()

	// 等待所有任务完成
	wg.Wait()
}
