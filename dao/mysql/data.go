package mysql

import (
	"blog/models"
	"errors"
)

const (
	essayNotExist = "没有找到该文章"
)

func GetDataAboutClassifyKind(data *[]string) (err error) {
	sqlStr := `SELECT ClassifyKindName FROM classifyKind `
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyDetails(data *[]models.DataAboutClassifyDetails) (err error) {
	sqlStr := `SELECT classifyKindName,classifyName,classifyRoute FROM classify`
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyEssayName(data *[]models.ClassifyIncludeEssay) (err error) {
	sqlStr := `SELECT essayName,essayKind,essayRoute FROM essay`
	err = db.Select(data, sqlStr)
	return
}

func UpdateEssay(data *models.UpdateEssay) (err error) {
	sqlStr := `UPDATE essay SET essayName= ?,essayKind = ? ,essayContent = ?,essayRoute = ? WHERE essayName = ?`
	result, err := db.Exec(sqlStr, data.EssayName, data.EssayKind, data.EssayContent, data.EssayRoute, data.EssayOldName)
	if err != nil {
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return
}
