package mysql

import (
	"blog/models"
	"fmt"
	"time"
)

func GetEssaySnowflakeID(id int) (eid int64, err error) {
	sqlStr := `SELECT eid FROM essay WHERE id = ?`
	if err = db.Get(&eid, sqlStr, id); err != nil {
		return 0, fmt.Errorf("no essay found with id %d", id)
	}
	return eid, nil
}

func GetDataAboutKind(data *[]models.DataAboutKind) error {
	sqlStr := `SELECT name,icon,id FROM kind `
	return db.Select(data, sqlStr)
}

func GetDataAboutClassifyDetails(data *[]models.DataAboutClassify) error {
	sqlStr := `SELECT kind,name,router,id FROM classify`
	return db.Select(data, sqlStr)
}

func GetDataAboutClassifyEssayMsg(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize

	// 使用 LIMIT 和 OFFSET 实现分页

	sqlStr1 := `SELECT name, kind, router, introduction, id, createdTime 
               FROM essay 
               WHERE name!='init' AND kind = ?
               ORDER BY id DESC 
               LIMIT ? OFFSET ?`
	sqlStr2 := `SELECT name, kind, router, introduction, id, createdTime 
               FROM essay 
               WHERE name!='init' 
               ORDER BY id DESC 
               LIMIT ? OFFSET ?`
	if query.Classify != "" {
		if err := db.Select(data.EssayList, sqlStr1, query.Classify, query.PageSize, offset); err != nil {
			return err
		}
	}
	if err := db.Select(data.EssayList, sqlStr2, query.PageSize, offset); err != nil {
		return err
	}
	// 计算总页数（向上取整）
	totalItems := len(*data.EssayList)
	totalPages := (totalItems + query.PageSize - 1) / query.PageSize
	//query.PageSize - 1) / query.PageSize 产生的结果加上任何不满分页值都会>=1 从而实现向上取整
	data.TotalPage = totalPages
	return nil
}

func GetEssayData(data *models.EssayData, id int) error {
	//在这里得到次数并添加
	sqlStr := `SELECT content,name,id,introduction,kind,createdTime,updatedTime,eid FROM essay where id = ?`
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
