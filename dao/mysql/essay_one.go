package mysql

import (
	"blog/models"
	"blog/pkg/snowflake"
	"blog/utils"
	"errors"
	"fmt"
)

func GetEssayData(data *models.EssayContent, id int) (err error) {
	sqlStr := `SELECT id, name,kind_id, content, introduction, createdTime, visitedTimes
			FROM essay where id = ?`

	if err = db.Get(data, sqlStr, id); err != nil {
		return err
	}

	if err = GetNearbyEssays(&data.NearEssayList, data.KindID, data.Id); err != nil {
		return err
	}
	return nil
}

func GetNearbyEssays(data *[]models.EssayData, kID int, eID int) error {
	sqlStr := `
		(SELECT id, name, kind_id, introduction, createdTime, imgUrl
			FROM essay 
			WHERE kind_id = ? AND id < ?
			ORDER BY id 
			LIMIT 2)
		UNION ALL
		(SELECT id, name,  kind_id, introduction, createdTime, imgUrl
			FROM essay 
			WHERE kind_id = ? AND id > ?
			ORDER BY id 
		LIMIT 2)
    `
	return db.Select(data, sqlStr, kID, eID, kID, eID)
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
