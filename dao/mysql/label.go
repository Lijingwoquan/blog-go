package mysql

import (
	"blog/models"
)

const (
	createLabelFailed = "创建label失败"
	updateLabelFailed = "更新label失败"
	deleteLabelFailed = "删除label失败"
)

func GetLabelList(data *[]models.LabelData) error {
	sqlStr := `SELECT name,id FROM label`
	return db.Select(data, sqlStr)
}

func CreateLabel(l *models.LabelParams) error {
	sqlStr := `INSERT INTO label (name) VALUES(:name)`
	result, err := db.Exec(sqlStr, l)
	if err != nil {
		return err
	}
	affect, err := result.RowsAffected()
	return noAffectedRowErr(affect, err, createLabelFailed)
}

func UpdateLabel(l *models.LabelParams) (err error) {

	return err
}

func DeleteLabel(l *models.LabelParams) (err error) {

	return err
}
