package mysql

import (
	"blog/models"
	"database/sql"
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

func GetAllDataAboutClassify(data *[]models.DataAboutClassify) error {
	sqlStr := `SELECT kind,name,router,id FROM classify`
	return db.Select(data, sqlStr)
}

func GetOneDataAboutClassify(data *models.DataAboutClassify) error {
	sqlStr := `SELECT kind,name,router,id FROM classify WHERE name = ?`
	return db.Get(data, sqlStr, data.Name)
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

func GetEssayData(data *models.EssayData, id int) (err error) {
	//在这里得到次数并添加
	sqlStr := `SELECT content,name,id,introduction,router,kind,createdTime,updatedTime,eid,imgUrl,advertiseMsg,advertiseImg,advertiseHref FROM essay where id = ?`

	if err = db.Get(data, sqlStr, id); err != nil {
		return err
	}
	if data.Last, data.Next, err = GetAdjacentEssay(id, data.Kind); err != nil {
		return err
	}
	return nil
}

func GetAdjacentEssay(currentID int, currentKind string) (models.AdjacentEssay, models.AdjacentEssay, error) {
	var lastEssay, nextEssay models.AdjacentEssay
	var lastID, nextID sql.NullInt64
	var lastName, nextName sql.NullString

	sqlStr := `
        SELECT 
            (SELECT id FROM essay WHERE id < ? AND kind = ? ORDER BY id DESC LIMIT 1) AS last_id,
            (SELECT name FROM essay WHERE id < ? AND kind = ? ORDER BY id DESC LIMIT 1) AS last_name,
            (SELECT id FROM essay WHERE id > ? AND kind = ? ORDER BY id ASC LIMIT 1) AS next_id,
            (SELECT name FROM essay WHERE id > ? AND kind = ? ORDER BY id ASC LIMIT 1) AS next_name
    `

	args := []interface{}{currentID, currentKind, currentID, currentKind, currentID, currentKind, currentID, currentKind}

	err := db.QueryRow(sqlStr, args...).Scan(&lastID, &lastName, &nextID, &nextName)
	if err != nil {
		return models.AdjacentEssay{}, models.AdjacentEssay{}, err
	}

	lastEssay = models.AdjacentEssay{Id: int(lastID.Int64), Name: lastName.String}
	nextEssay = models.AdjacentEssay{Id: int(nextID.Int64), Name: nextName.String}

	return lastEssay, nextEssay, nil
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
