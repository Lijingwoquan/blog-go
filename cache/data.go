package cache

import (
	"blog/help"
	"blog/models"
	"go.uber.org/zap"
	"sync"
)

var (
	GlobalDataAboutIndex = models.DataAboutIndex{}
	Mu                   sync.RWMutex // 读写锁
)

func UpdateDataAboutIndex() {
	errCh := make(chan error)
	done := make(chan bool)
	go func() {
		if err := getDataAboutIndex(); err != nil {
			errCh <- err
		}
		done <- true
	}()
	go func() {
		for err := range errCh {
			zap.L().Error("happen err in cache Update:%v", zap.Error(err))
		}
	}()
	go func() {
		<-done
		close(errCh)
	}()
}

func getDataAboutIndex() error {
	Mu.Lock()
	defer Mu.Unlock()

	if err := help.ResponseDataAboutIndex(&GlobalDataAboutIndex); err != nil {
		return err
	}
	return nil
}
