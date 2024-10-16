package mysql

import (
	"blog/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func GetEssaySnowflakeID(id int) (eid int64, err error) {
	sqlStr := `SELECT eid FROM essay WHERE id = ?`
	if err = db.Get(&eid, sqlStr, id); err != nil {
		return 0, fmt.Errorf("no essay found with id %d", id)
	}
	return eid, nil
}

func GetKindList(data *[]models.DataAboutKind) error {
	sqlStr := `SELECT name,icon,id,router,essayCount FROM kind `
	return db.Select(data, sqlStr)
}

func GetLabelList(data *[]models.DataAboutLabel) error {
	sqlStr := `SELECT kind,name,router,id FROM label`
	return db.Select(data, sqlStr)
}

func GetRecommendEssayList(data *[]models.DataAboutEssay) error {
	sqlStr := `SELECT id, name, createdTime, imgUrl FROM essay WHERE ifRecommend = true`
	return db.Select(data, sqlStr)
}

func GetEssayList(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize
	baseSelect := `SELECT id, name, label, kind, ifRecommend,introduction, createdTime, visitedTimes, eid, imgUrl
	          FROM essay`
	baseCount := `SELECT COUNT(*) FROM essay`
	where := make([]string, 0)
	args := make([]interface{}, 0)
	if query.Label != "" {
		where = append(where, "label = ?")
		args = append(args, query.Label)
	}
	if query.Kind != "" {
		where = append(where, "kind = ?")
		args = append(args, query.Kind)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}
	sqlStr := fmt.Sprintf("%s %s ORDER BY id DESC LIMIT ? OFFSET ?", baseSelect, whereClause)
	countSqlStr := fmt.Sprintf("%s %s", baseCount, whereClause)
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
	totalPages := totalCount/query.PageSize + (+query.PageSize-1)/query.PageSize
	// 设置总页数
	data.TotalPages = totalPages

	return nil
}

func GetOneDataAboutClassify(data *models.DataAboutLabel) error {
	sqlStr := `SELECT kind,name,router,id FROM label WHERE name = ?`
	return db.Get(data, sqlStr, data.Name)
}

func GetAllEssay(data *[]models.DataAboutEssay) error {
	sqlStr := `SELECT  name, kind, introduction, id
               FROM essay 
               ORDER BY id DESC`
	return db.Select(data, sqlStr)
}

func GetEssayData(data *models.EssayData, id int) (err error) {
	sqlStr := `SELECT content,name,id,introduction,kind,createdTime,updatedTime,eid,imgUrl FROM essay where id = ?`

	var advertiseMsg sql.NullString
	var advertiseImg sql.NullString
	var advertiseHref sql.NullString

	if err = db.QueryRow(sqlStr, id).Scan(&data.Content, &data.Name, &data.Id, &data.Introduction, &data.Router, &data.Kind, &data.CreatedTime, &data.UpdatedTime, &data.Eid, &data.ImgUrl, &advertiseMsg, &advertiseImg, &advertiseHref); err != nil {
		return err
	}

	data.AdvertiseMsg = advertiseMsg.String
	data.AdvertiseImg = advertiseImg.String
	data.AdvertiseHref = advertiseHref.String

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
            (SELECT id FROM essay WHERE id > ? AND kind = ? ORDER BY id  LIMIT 1) AS next_id,
            (SELECT name FROM essay WHERE id > ? AND kind = ? ORDER BY id  LIMIT 1) AS next_name
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
