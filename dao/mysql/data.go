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
	var sqlStr string
	var countSqlStr string
	var args []interface{}
	if query.Classify != "" {
		// 根据分类查询
		sqlStr = `SELECT name, kind, router, introduction, id, createdTime,eid,imgUrl
               FROM essay 
               WHERE  kind = ?
               ORDER BY id DESC 
               LIMIT ? OFFSET ?`
		countSqlStr = `SELECT COUNT(*) 
               FROM essay 
               WHERE kind = ?`
		args = append(args, query.Classify)
	} else {
		//返回首页文章列表
		sqlStr = `SELECT name, kind, router, introduction, id, createdTime ,eid,imgUrl
               FROM essay 
               ORDER BY id DESC 
               LIMIT ? OFFSET ?`
		countSqlStr = `SELECT COUNT(*)  
               FROM essay`
	}
	args = append(args, query.PageSize, offset)

	//执行分页查询
	if err := db.Select(data.EssayList, sqlStr, args...); err != nil {
		return err
	}
	//执行总记录统计
	var totalCount int
	if err := db.Get(&totalCount, countSqlStr, args[:len(args)-2]...); err != nil {
		return err
	}

	//计算总记录数
	totalPages := (totalCount + query.PageSize - 1) / query.PageSize

	// 设置总页数
	data.TotalPages = totalPages
	return nil
}

func GetAllEssay(data *[]models.DataAboutEssay) error {
	sqlStr := `SELECT  name, kind, router, introduction, id 
               FROM essay 
               ORDER BY id DESC`
	return db.Select(data, sqlStr)
}

func GetEssayData(data *models.EssayData, id int) error {
	//在这里得到次数并添加
	sqlStr := `SELECT content,name,id,introduction,router,kind,createdTime,updatedTime,eid FROM essay where id = ?`
	return db.Get(data, sqlStr, id)
}

func CleanupInvalidTokens() error {
	now := time.Now()
	sqlStr := `DELETE FROM tokenInvalid WHERE expiration < ? `
	_, err := db.Exec(sqlStr, now)
	return err
}

func SaveVisitedTimes(visitedTimesChangedMap map[int64]int64) error {
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
