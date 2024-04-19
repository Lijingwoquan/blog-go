package mysql

import (
	"blog/models"
	"time"
)

func GetDataAboutKind(data *[]models.DataAboutKind) error {
	sqlStr := `SELECT name,icon,id FROM kind `
	return db.Select(data, sqlStr)
}

func GetDataAboutClassifyDetails(data *[]models.DataAboutClassify) error {
	sqlStr := `SELECT kind,name,router,id FROM classify`
	return db.Select(data, sqlStr)
}

func GetDataAboutClassifyEssayMsg(data *[]models.DataAboutEssay) error {
	sqlStr := `SELECT name,kind,router,introduction,id,createdTime FROM essay WHERE name!='init'  ORDER BY id DESC`
	return db.Select(data, sqlStr)
}

func GetEssayData(data *models.EssayData, id int) error {
	//在这里得到次数并添加
	sqlStr := `SELECT content,name,eid,introduction,kind,createdTime,updatedTime FROM essay where id = ?`
	return db.Get(data, sqlStr, id)
}

func CleanupInvalidTokens() error {
	now := time.Now()
	sqlStr := `DELETE FROM tokenInvalid WHERE expiration < ? `
	_, err := db.Exec(sqlStr, now)
	return err
}

func SaveVisitedTimes(visitedTimesChangedMap map[int64]int64) (err error) {
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

	sqlStr := `UPDATE essay SET visitedTimes = ? WHERE eid = ?`

	for eid, vt := range visitedTimesChangedMap {
		if _, err := tx.Exec(sqlStr, vt, eid); err != nil {
			return err
		}
	}
	return nil
}

func GetVisitedTimesFromMySQL(eid string) (vt int64, err error) {
	sqlStr := `SELECT visitedTimes FROM essay WHERE  eid = ?`
	if err = db.Get(&vt, sqlStr, eid); err != nil {
		return 0, err
	}
	return vt, nil
}
