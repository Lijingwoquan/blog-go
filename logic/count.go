package logic

import (
	"blog/dao/redis"
	"blog/models"
)

func GetUserIpCount(setKind *models.SetKind) (err error) {
	return redis.GetUserIpCount(setKind)
}
