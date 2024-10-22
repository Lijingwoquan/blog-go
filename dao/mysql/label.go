package mysql

import (
	"blog/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func GetLabelList(data *[]models.LabelData) error {
	sqlStr := `SELECT name,id FROM label`
	return db.Select(data, sqlStr)
}

func CreateLabel(l *models.LabelParams) error {
	sqlStr := `INSERT INTO label (name) VALUES(:name)`
	result, err := db.NamedExec(sqlStr, l)
	if err != nil {
		return err
	}
	affect, err := result.RowsAffected()
	return noAffectedRowErr(affect, err, "create label failed")
}

func DeleteLabel(id int) error {
	return withTx(func(tx *sqlx.Tx) error {
		if err := deleteLabelInEssayLabel(tx, id); err != nil {
			return fmt.Errorf("deleteLabelFromEssayLabel failed,err:%w", err)
		}
		if err := deleteLabelInLabel(tx, id); err != nil {
			return fmt.Errorf("deleteLabelFromLabel failed,err:%w", err)
		}
		return nil
	})
}

func deleteLabelInLabel(tx *sqlx.Tx, id int) error {
	// 删除label
	sqlStr := `DELETE FROM label WHERE id = ?`
	ret, err := tx.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	affect, err := ret.RowsAffected()
	fmt.Println(err)
	return noAffectedRowErr(affect, err, "have not change any row")
}

func deleteLabelInEssayLabel(tx *sqlx.Tx, id int) error {
	sqlStr := `DELETE FROM essay_label WHERE label_id = ?`
	ret, err := tx.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	affect, err := ret.RowsAffected()
	return noAffectedRowErr(affect, err, "delete label within essay_label failed,err")
}

func UpdateLabel(l *models.LabelUpdateParams) error {

	return nil
}
