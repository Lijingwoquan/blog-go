package mysql

import "blog/models"

func GetKindList(data *[]models.KindData) error {
	sqlStr := `SELECT name,icon,id,essayCount FROM kind `
	return db.Select(data, sqlStr)
}

func CheckKind(id int) (oldName string, err error) {
	sqlStr := `SELECT name FROM kind WHERE  id = ?`
	err = db.QueryRow(sqlStr, id).Scan(&oldName)
	return
}

func UpdateKind(oldName string, k *models.UpdateKindParams) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	//更新classify的name
	sqlStr1 := `UPDATE kind SET name =?,icon=? WHERE  id = ?`
	if _, err = tx.Exec(sqlStr1, k.Name, k.Icon, k.ID); err != nil {
		return err
	}
	//更新符合classify的kind
	sqlStr2 := `UPDATE  classify SET kind = ? WHERE kind = ?`
	_, err = tx.Exec(sqlStr2, k.Name, oldName)
	return err
}
