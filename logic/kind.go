package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

// UpdateKind 更新总纲逻辑
func UpdateKind(k *models.KindParams) (err error) {
	return mysql.UpdateKind(k)
}
