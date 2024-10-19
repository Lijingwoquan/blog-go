package cache

import (
	"blog/help"
	"blog/models"
	"go.uber.org/zap"
)

var (
	globalDataAboutEssayList = new([]models.DataAboutEssay)
)

func GetEssayListInit() (*[]models.DataAboutEssay, error) {
	if err := help.ResponseDataAboutEssayList(globalDataAboutEssayList); err != nil {
		zap.L().Error("help.ResponseDataAboutEssayList(globalDataAboutEssayList) filed,err:", zap.Error(err))
		return nil, err
	}
	return globalDataAboutEssayList, nil
}

func GetAllEssayList() *[]models.DataAboutEssay {
	return globalDataAboutEssayList
}

func UpdateDataAboutEssayList() {
	errCh := make(chan error)
	done := make(chan bool)
	go func() {
		if _, err := GetEssayListInit(); err != nil {
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
