package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

func CreateLabel(l *models.LabelParams) (err error) {
	return mysql.CreateLabel(l)
}

func UpdateLabel(l *models.LabelParams) (err error) {
	//2.由传进来的id查询数据库 进行更新
	return mysql.UpdateLabel(l)
}

func DeleteLabel(l *models.LabelParams) (err error) {
	//2.由传进来的id查询数据库 进行更新
	return mysql.DeleteLabel(l)
}
