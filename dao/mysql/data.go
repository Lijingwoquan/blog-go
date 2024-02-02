package mysql

import "blog/models"

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

//先得到
