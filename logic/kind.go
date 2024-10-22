package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

func CreateKind(k *models.KindParams) (err error) {
	return mysql.CreateKind(k)
}

func DeleteKind(k *models.KindParams) (err error) {

	return err
}

func UpdateKind(k *models.KindUpdateParams) (err error) {
	return mysql.UpdateKind(k)
}
