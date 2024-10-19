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
	errCh := make(chan error, tickerTaskCount)

	tasks := []func() error{
		func() error {
			err := InitIndexData()
			return err
		},
		func() error {
			_, err := GetEssayListInit()
			return err
		},
		getMaliciousMap,
	}

	for _, task := range tasks {
		wg.Add(1)
		go func(t func() error) {
			defer wg.Done()
			if err := t(); err != nil {
				errCh <- err
			}
		}(task)
	}

	go func() {
		wg.Wait()
		close(errCh) //关闭通道避免死循环
	}()

	for err := range errCh {
		zap.L().Error("Error in cache Init", zap.Error(err))
	}
}
