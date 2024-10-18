package logic

import (
	"blog/cache"
	"blog/models"
)

func GetIndexData() (data *models.DataAboutIndex, err error) {
	// 从缓存中拿到数据
	if data, err = cache.GetIndexData(); err != nil {
		return nil, err
	}
	return data, nil
}
