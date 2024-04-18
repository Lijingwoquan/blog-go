package cache

import (
	"blog/logic"
	"blog/models"
	"go.uber.org/zap"
	"sync"
)

var (
	GlobalDataAboutIndex = models.DataAboutIndex{}
	Mu                   sync.RWMutex // 读写锁
)

func Init() {
	errCh := make(chan error)
	done := make(chan bool)
	//updateDataAboutIndex
	go func() {
		if err := updateDataAboutIndex(); err != nil {
			errCh <- err
		}
		done <- true
	}()
	//错误处理
	go func() {
		for err := range errCh {
			zap.L().Error("happen err in cache Init:%v", zap.Error(err))
		}
	}()
	//关闭通道
	go func() {
		<-done
		close(errCh)
	}()
}

func Update() {
	errCh := make(chan error)
	done := make(chan bool)
	go func() {
		if err := updateDataAboutIndex(); err != nil {
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

// 更新数据
func updateDataAboutIndex() error {
	Mu.Lock()
	defer Mu.Unlock()

	if err := logic.ResponseDataAboutIndex(&GlobalDataAboutIndex); err != nil {
		return err
	}
	return nil
}
