package cache

import (
	"blog/help"
	"blog/models"
	"go.uber.org/zap"
)

var (
	globalDataAboutIndex     = new(models.DataAboutIndex)
	globalDataAboutEssayList = new([]models.DataAboutEssay)
	Error                    error
)

func UpdateDataAboutIndex() {
	errCh := make(chan error)
	done := make(chan bool)
	go func() {
		if _, err := GetDataAboutIndex(); err != nil {
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

func UpdateDataAboutEssayList() {
	errCh := make(chan error)
	done := make(chan bool)
	go func() {
		if _, err := GetEssayList(); err != nil {
			errCh <- err
		}
		done <- true
	}()
	go func() {
		for err := range errCh {
			zap.L().Error("happen err in cache UpdateDataAboutEssayList:%v", zap.Error(err))
		}
	}()
	go func() {
		<-done
		close(errCh)
	}()
}

func GetDataAboutIndex() (*models.DataAboutIndex, error) {
	if Error = help.ResponseDataAboutIndex(globalDataAboutIndex); Error != nil {
		return nil, Error
	}
	return globalDataAboutIndex, nil
}

func GetEssayList() (*[]models.DataAboutEssay, error) {
	if Error = help.ResponseDataAboutEssayList(globalDataAboutEssayList); Error != nil {
		return nil, Error
	}
	return globalDataAboutEssayList, nil
}
