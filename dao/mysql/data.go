package mysql

import (
	"blog/models"
	"time"
)

func GetDataAboutKind(data *[]models.DataAboutKind) (err error) {
	sqlStr := `SELECT name,icon,id FROM kind `
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyDetails(data *[]models.DataAboutClassify) (err error) {
	sqlStr := `SELECT kind,name,router,id FROM classify`
	err = db.Select(data, sqlStr)
	return
}

func GetDataAboutClassifyEssayMsg(data *[]models.DataAboutEssay) (err error) {
	sqlStr := `SELECT name,kind,router,introduction,id,createdTime FROM essay WHERE name!='init'  ORDER BY id DESC`
	err = db.Select(data, sqlStr)
	return
}

func GetEssayData(data *models.EssayData, id int) (err error) {
	sqlStr := `SELECT content,name,introduction,kind,createdTime,updatedTime FROM essay where id = ?`
	err = db.Get(data, sqlStr, id)
	return
}

func CleanupInvalidTokens() (err error) {
	now := time.Now()
	sqlStr := `DELETE FROM tokenInvalid WHERE expiration < ? `
	_, err = db.Exec(sqlStr, now)
	return err
}
