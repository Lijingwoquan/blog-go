package mysql

import (
	"blog/models"
	"blog/utils"
	"errors"
	"fmt"
)

const (
	invalidLabelIds = "labels参撒无效"
)

func GetEssayData(data *models.EssayContent, id int) (err error) {
	sqlStr := `
		SELECT e.id,e.name,e.kind_id, e.content, e.introduction, e.created_time, e.visited_times,
			k.name AS kind_name
		FROM essay e 
		LEFT JOIN kind k on e.kind_id = k.id
		where e.id = ?
		`

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
		(SELECT e.id, e.name, e.kind_id, e.introduction, e.created_time, e.img_url,
		 	k.name AS kind_name
			FROM essay e 
			LEFT JOIN kind k on k.id = e.kind_id
			WHERE e.kind_id = ? AND e.id < ?
			ORDER BY e.id 
			LIMIT 2)
		UNION ALL
		(SELECT e.id, e.name, e.kind_id, e.introduction, e.created_time, e.img_url,
		 	k.name AS kindName
			FROM essay e 
			LEFT JOIN kind k on k.id = e.kind_id
			WHERE e.kind_id = ? AND e.id > ?
			ORDER BY e.id 
		LIMIT 2)
    `
	return db.Select(data, sqlStr, kID, eID, kID, eID)
}

func CreateEssay(e *models.EssayParams) (err error) {
	if e.CreatedTime, err = utils.GetChineseTime(); err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 在essay表添加数据
	sqlStr := `
		INSERT INTO essay(name, kind_id, if_top,content, if_recommend, introduction, created_time, img_url) 
			values(:name, :kind_id, :if_top,:content, :if_recommend, :introduction, :created_time, :img_url)`

	result, err := tx.NamedExec(sqlStr, e)
	if err != nil {
		return err
	}

	essayID, _ := result.LastInsertId()
	if err != nil {
		return err
	}

	//  添加标签关联
	if len(e.LabelIds) > 0 {
		sqlStr2 := `
            INSERT INTO essay_label (essay_id, label_id, label_name) 
            SELECT ?, ?, name
            FROM label 
            WHERE id = ?
        `
		for _, labelID := range e.LabelIds {
			result, err := tx.Exec(sqlStr2, essayID, labelID, labelID)
			if err != nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err = noAffectedRowErr(affected, err, invalidLabelIds); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("请求参数缺少label")
	}

	// 在相应的kind表中增加essay_count的值
	sqlStr3 := `
		UPDATE kind SET essay_count = essay_count + 1
		WHERE id = ?`

	if _, err = tx.Exec(sqlStr3, e.KindID); err != nil {
		return err
	}

	return tx.Commit()
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
