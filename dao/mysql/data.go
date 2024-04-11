package mysql

import (
	"blog/dao/redis"
	"blog/models"
	"fmt"
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
	//在这里得到次数并添加
	sqlStr := `SELECT content,name,eid,introduction,kind,createdTime,updatedTime FROM essay where id = ?`
	err = db.Get(data, sqlStr, id)
	if err != nil {
		return err
	}
	data.VisitedTimes, err = redis.GetVisitedTimes(fmt.Sprintf("%d", data.Eid))
	return
}

func CleanupInvalidTokens() (err error) {
	now := time.Now()
	sqlStr := `DELETE FROM tokenInvalid WHERE expiration < ? `
	_, err = db.Exec(sqlStr, now)
	return err
}

func SaveVisitedTimes(visitedTimesChangedMap *map[int64]int64) (err error) {
	for eid, vt := range *visitedTimesChangedMap {
		sqlStr := `UPDATE essay SET visitedTimes = ? WHERE eid = ?`
		_, err := db.Exec(sqlStr, vt, eid)
		if err != nil {
			return fmt.Errorf("_, err := db.Exec(sqlStr, vt, eid) in SaveVisitedTimes failed,err:%v", err)
		}
	}
	return nil
}
