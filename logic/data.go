package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

func GetEssayData(essay *models.EssayContent, id int) (err error) {
	//从数据库查询文章对应id的内容
	err = mysql.GetEssayData(essay, id)
	return
}
