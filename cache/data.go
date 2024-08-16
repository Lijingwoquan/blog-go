package cache

import (
	"blog/help"
	"blog/models"
	"go.uber.org/zap"
)

var (
	GlobalDataAboutIndex = models.DataAboutIndex{}
	Error                error
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
	if Error = help.ResponseDataAboutIndex(&GlobalDataAboutIndex); Error != nil {
		return Error
	}
	return nil
}
