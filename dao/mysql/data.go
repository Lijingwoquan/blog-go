package mysql

import (
	"blog/models"
)

const (
	essayNotExist = "没有找到该文章"
)

func GetDataAboutClassifyKind(data *[]models.DataAboutClassify) (err error) {
	sqlStr := `SELECT name,icon FROM classifyKind `
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyDetails(data *[]models.DataAboutClassifyDetails) (err error) {
	sqlStr := `SELECT kind,name,router,id FROM classify`
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyEssayMsg(data *[]models.DataAboutEssay) (err error) {
	sqlStr := `SELECT name,kind,router,introduction,id FROM essay`
	err = db.Select(data, sqlStr)
	return
}
