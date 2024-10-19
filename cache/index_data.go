package cache

import (
	"blog/help"
	"blog/models"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
)

var (
	globalDataAboutIndex = new(models.IndexData)
)

func InitIndexData() error {
	if err := help.ResponseIndexData(globalDataAboutIndex); err != nil {
		zap.L().Error("help.ResponseDataAboutIndex(globalDataAboutIndex) failed,err:", zap.Error(err))
		return err
	}
	return nil
}

func GetIndexData(data **models.IndexData) {
	for index := range (*globalDataAboutIndex).LabelList {
		(*globalDataAboutIndex).LabelList[index].Color = fmt.Sprintf("rgb(%d,%d,%d)", rand.Intn(256), rand.Intn(256), rand.Intn(256))
	}
	*data = globalDataAboutIndex
}

func UpdateIndexData() {
	errCh := make(chan error)
	done := make(chan bool)
	go func() {
		if err := InitIndexData(); err != nil {
			errCh <- err
		}
		done <- true
	}()
	go func() {
		for err := range errCh {
			zap.L().Error("happen err in cache UpdateDataAboutIndex:%v", zap.Error(err))
		}
	}()
	go func() {
		<-done
		close(errCh)
	}()
}
