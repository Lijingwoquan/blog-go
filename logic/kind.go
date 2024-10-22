package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

func CreateKind(k *models.KindParams) error {
	return mysql.CreateKind(k)
}

func DeleteKind(id int) error {
	return mysql.DeleteKind(id)
}

func UpdateKind(k *models.KindUpdateParams) error {
	return mysql.UpdateKind(k)
}
