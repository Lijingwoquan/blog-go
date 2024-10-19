package logic

import (
	"blog/cache"
	"blog/models"
)

func GetIndexData(data **models.IndexData) {
	// 从缓存中拿到数据
	cache.GetIndexData(data)
}
