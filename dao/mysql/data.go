package mysql

import (
	"blog/models"
	"database/sql"
	"fmt"
	"strconv"
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
	sqlStr := `SELECT name,icon,id,essayCount FROM kind `
	return db.Select(data, sqlStr)
}

func GetLabelList(data *[]models.DataAboutLabel) error {
	sqlStr := `SELECT name,id FROM label`
	return db.Select(data, sqlStr)
}

func GetRecommendEssayList(data *[]models.DataAboutEssay) error {
	sqlStr := `
		SELECT e.id, e.name, e.createdTime, e.imgUrl 
		FROM essay e 
		WHERE ifRecommend = true
	`
	return db.Select(data, sqlStr)
}

func test(data *[]models.DataAboutEssay) error {
	type rawData struct {
		models.DataAboutEssay
		LabelIDs   string `db:"label_ids"`
		LabelNames string `db:"label_names"`
	}
	var rawDataList = new([]rawData)
	sqlStr := `
		SELECT e.id, e.name, e.createdTime, e.imgUrl,e.kind_id,
		       k.name AS kind_name,
		       GROUP_CONCAT(el.label_id) AS label_ids ,GROUP_CONCAT(el.label_name) AS label_names
		FROM essay e
		LEFT JOIN kind k on e.kind_id = k.id
		LEFT JOIN essay_label el on e.id = el.essay_id
		WHERE  e.ifRecommend = true
		GROUP BY e.id, e.name, e.createdTime, e.imgUrl, e.kind_id,  k.name
	`
	var err error
	if err = db.Select(rawDataList, sqlStr); err != nil {
		return err
	}
	*data = make([]models.DataAboutEssay, len(*rawDataList))
	for i, raw := range *rawDataList {
		(*data)[i] = raw.DataAboutEssay
		if raw.LabelNames != "" && raw.LabelIDs != "" {
			ids := strings.Split(raw.LabelIDs, ",")
			names := strings.Split(raw.LabelNames, ",")
			(*data)[i].LabelList = make([]models.Label, len(ids))
			for j := range ids {
				id, _ := strconv.Atoi(ids[j])
				(*data)[i].LabelList[j] = models.Label{
					ID:   id,
					Name: names[j],
				}
			}
		}
	}
	return err
}

func GetEssayList(data *models.DataAboutEssayListAndPage, query models.EssayQuery) error {
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize
	baseSelect := `
		 SELECT  e.id, e.name, e.kind_id, e.ifRecommend, e.introduction, e.createdTime, e.visitedTimes, e.imgUrl,
            k.name AS kind_name,
            GROUP_CONCAT(el.label_id) AS label_ids,
            GROUP_CONCAT(el.label_name) AS label_names
        FROM essay e
        LEFT JOIN kind k ON e.kind_id = k.id
        LEFT JOIN essay_label el ON e.id = el.essay_id`

	baseCount := `
		 SELECT COUNT(DISTINCT e.id)
        FROM essay e 
        LEFT JOIN kind k ON e.kind_id = k.id
        LEFT JOIN essay_label el ON e.id = el.essay_id
		`

	where := make([]string, 0)
	args := make([]interface{}, 0)
	if query.LabelID != 0 {
		where = append(where, "el.label_id = ?")
		args = append(args, query.LabelID)
	}
	if query.KindID != 0 {
		where = append(where, "e.kind_id = ?")
		args = append(args, query.KindID)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// 添加GROUP BY子句
	groupBy := "GROUP BY e.id, e.name, e.kind_id, e.ifRecommend, e.introduction, e.createdTime, e.visitedTimes, e.imgUrl, k.name"

	// 构建完整的SQL语句
	sqlStr := fmt.Sprintf("%s %s %s ORDER BY e.id DESC LIMIT ? OFFSET ?",
		baseSelect, whereClause, groupBy)
	countSqlStr := fmt.Sprintf("%s %s", baseCount, whereClause)
	args = append(args, query.PageSize, offset)
	// 创建临时结构体来接收原始数据
	type rawData struct {
		models.DataAboutEssay
		LabelIDs   string `db:"label_ids"`
		LabelNames string `db:"label_names"`
	}

	var rawDataList []rawData

	// 执行分页查询
	if err := db.Select(&rawDataList, sqlStr, args...); err != nil {
		return err
	}

	// 处理查询结果
	data.EssayList = make([]models.DataAboutEssay, len(rawDataList))
	for i, raw := range rawDataList {
		data.EssayList[i] = raw.DataAboutEssay
		// 处理标签数据
		if raw.LabelIDs != "" && raw.LabelNames != "" {
			ids := strings.Split(raw.LabelIDs, ",")
			names := strings.Split(raw.LabelNames, ",")
			data.EssayList[i].LabelList = make([]models.Label, len(ids))

			for j := range ids {
				id, _ := strconv.Atoi(ids[j])
				name := names[j]
				data.EssayList[i].LabelList[j] = models.Label{
					ID:   id,
					Name: name,
				}
			}
		}
	}

	// 执行总记录统计
	var totalCount int
	if err := db.Get(&totalCount, countSqlStr, args[:len(args)-2]...); err != nil {
		return err
	}

	// 计算总页数
	data.TotalPages = (totalCount + query.PageSize - 1) / query.PageSize

	return nil
}

func GetOneDataAboutClassify(data *models.DataAboutLabel) error {
	sqlStr := `SELECT name,id FROM label WHERE name = ?`
	return db.Get(data, sqlStr, data.Name)
}

func GetAllEssay(data *[]models.DataAboutEssay) error {
	sqlStr := `SELECT  name,  introduction, id
               FROM essay 
               ORDER BY id DESC`
	return db.Select(data, sqlStr)
}

func GetEssayData(data *models.EssayData, id int) (err error) {
	sqlStr := `SELECT content,name,id,introduction,kind,createdTime,updatedTime,eid,imgUrl FROM essay where id = ?`

	var advertiseMsg sql.NullString
	var advertiseImg sql.NullString
	var advertiseHref sql.NullString

	if err = db.QueryRow(sqlStr, id).Scan(&data.Content, &data.Name, &data.Id, &data.Introduction, &data.Kind, &data.CreatedTime, &data.Eid, &data.ImgUrl, &advertiseMsg, &advertiseImg, &advertiseHref); err != nil {
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
