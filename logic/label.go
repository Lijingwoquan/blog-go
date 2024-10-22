package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

func CreateLabel(l *models.LabelParams) (err error) {
	return mysql.CreateLabel(l)
}

func DeleteLabel(id int) (err error) {
	return mysql.DeleteLabel(id)
}

func UpdateLabel(l *models.LabelUpdateParams) (err error) {
	return mysql.UpdateLabel(l)
}
