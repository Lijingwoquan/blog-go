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

func GetDataAboutIndex() (*models.DataAboutIndex, error) {
	if Error = help.ResponseDataAboutIndex(globalDataAboutIndex); Error != nil {
		return nil, Error
	}
	return globalDataAboutIndex, nil
}

<<<<<<< HEAD
func GetEssayList() (*[]models.DataAboutEssay, error) {
	if Error = help.ResponseDataAboutEssayList(globalDataAboutEssayList); Error != nil {
		return nil, Error
=======
func GetEssayListInit() (*[]models.DataAboutEssay, error) {
	if err := help.ResponseDataAboutEssayList(globalDataAboutEssayList); err != nil {
		zap.L().Error("help.ResponseDataAboutEssayList(globalDataAboutEssayList) filed,err:", zap.Error(err))
		return nil, err
>>>>>>> dev
	}
	return globalDataAboutEssayList, nil
}

func GetAllEssayList() *[]models.DataAboutEssay {
	return globalDataAboutEssayList
}
