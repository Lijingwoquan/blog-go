package mysql

import "blog/models"

func GetKindList(data *[]models.KindData) error {
	sqlStr := `SELECT name,icon,id,essay_count FROM kind `
	return db.Select(data, sqlStr)
}

func CreateKind(k *models.KindParams) error {
	sqlStr := `INSERT INTO kind(name, icon) VALUES (:name,:icon)`
	_, err := db.NamedExec(sqlStr, k)
	return err
}

func DeleteKind(k *models.KindParams) (err error) {

	return err
}

func UpdateKind(k *models.KindUpdateParams) (err error) {

	return err
}
