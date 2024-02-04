package mysql

import (
	"blog/models"
)

const (
	essayNotExist = "没有找到该文章"
)

func GetDataAboutClassifyKind(data *[]string) (err error) {
	sqlStr := `SELECT name FROM classifyKind `
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyDetails(data *[]models.DataAboutClassifyDetails) (err error) {
	sqlStr := `SELECT kind,name,router FROM classify`
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyEssayName(data *[]models.ClassifyIncludeEssay) (err error) {
	sqlStr := `SELECT name,kind,router FROM essay`
	err = db.Select(data, sqlStr)
	return
}
