package mysql

import (
	"blog/models"
	"blog/pkg/snowflake"
	"blog/utils"
	"database/sql"
	"errors"
	"fmt"
)

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

func DeleteEssay(id int) error {
	sqlStr := `DELETE  FROM essay WHERE id=?`
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	var RowsAffected int64
	if RowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if RowsAffected == 0 {
		return errors.New(essayNotExist)
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

func GetEssaySnowflakeID(id int) (eid int64, err error) {
	sqlStr := `SELECT eid FROM essay WHERE id = ?`
	if err = db.Get(&eid, sqlStr, id); err != nil {
		return 0, fmt.Errorf("no essay found with id %d", id)
	}
	return eid, nil
}

// CheckEssayExist 检测文章是否存在
func CheckEssayExist(c *models.EssayParams) error {
	var err error
	sqlStr := `SELECT COUNT(*) FROM essay WHERE  kind = ? AND  name = ? `
	var count int
	if err = db.Get(&count, sqlStr, c.Kind, c.Name); err != nil {
		return err
	}

	if count > 0 {
		return errors.New(essayExist)
	}
	return nil
}

// CreateEssay 添加新文章
func CreateEssay(e *models.EssayParams) (erd int64, err error) {
	var formattedTime string
	if formattedTime, err = utils.GetChineseTime(); err != nil {
		return 0, err
	}
	eid := snowflake.GenID()
	sqlStr := `INSERT INTO essay(kind,name,content,Introduction,createdTime,updatedTime,eid,imgUrl) values(?,?,?,?,?,?,?,?)`
	if _, err = db.Exec(sqlStr, e.Kind, e.Name, e.Content, e.Introduction, formattedTime, formattedTime, eid, e.ImgUrl); err != nil {
		return 0, err
	}
	return eid, err
}

// UpdateEssayMsg 更新文章基本信息
func UpdateEssayMsg(data *models.UpdateEssayMsgParams) error {
	var err error
	var formattedTime string
	if formattedTime, err = utils.GetChineseTime(); err != nil {
		return err
	}
	sqlStr := `UPDATE essay SET name= ?,kind = ? ,content = ?,introduction=?,updatedTime=?,imgUrl=?,advertiseMsg=?,advertiseImg=?,advertiseHref = ? WHERE id = ?`
	result, err := db.Exec(
		sqlStr,
		data.Name, data.Kind, data.Content, data.Introduction, formattedTime, data.ImgUrl, data.AdvertiseMsg, data.AdvertiseImg, data.AdvertiseHref,
		data.Id)
	if err != nil {
		return err
	}
	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return nil
}
