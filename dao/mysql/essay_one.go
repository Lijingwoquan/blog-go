package mysql

import (
	"blog/models"
	"blog/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	invalidLabelIds = "labels参数无效"
	essayNoExist    = "该文章不存在"
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

// CreateEssay 主函数
func CreateEssay(e *models.EssayParams) error {
	var err error

	if len(e.LabelIds) == 0 {
		return fmt.Errorf(invalidLabelIds)
	}

	if e.CreatedTime, err = utils.GetChineseTime(); err != nil {
		return fmt.Errorf("get chinese time failed: %w", err)
	}

	return withTx(func(tx *sqlx.Tx) error {
		result, err := insertEssay(tx, e)
		if err != nil {
			return fmt.Errorf("insert essay failed: %w", err)
		}

		eid64, err := result.LastInsertId()
		if err != nil {
			return err
		}

		eid := int(eid64)

		if err := insertEssayLabels(tx, eid, e.LabelIds); err != nil {
			return fmt.Errorf("insert essay label failed: %w", err)
		}

		if err := updateKindEssayCount(tx, e.KindID, true); err != nil {
			return fmt.Errorf("update kind essay_count failed: %w", err)
		}
		return nil
	})
}

// 插入文章
func insertEssay(tx *sqlx.Tx, e *models.EssayParams) (sql.Result, error) {
	sqlStr := `
        INSERT INTO essay(name, kind_id, if_top, content, if_recommend, introduction, created_time, img_url) 
        VALUES (:name, :kind_id, :if_top, :content, :if_recommend, :introduction, :created_time, :img_url)
    `
	result, err := tx.NamedExec(sqlStr, e)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 插入文章标签关联
func insertEssayLabels(tx *sqlx.Tx, eid int, lIDs []int) error {
	sqlStr := `
        INSERT INTO essay_label (essay_id, label_id, label_name) 
        SELECT ?, ?, name
        FROM label 
        WHERE id = ?
    `

	for _, labelID := range lIDs {
		result, err := tx.Exec(sqlStr, eid, labelID, labelID)
		if err != nil {
			return err
		}

		affected, err := result.RowsAffected()
		if err = noAffectedRowErr(affected, err, invalidLabelIds); err != nil {
			return err
		}
	}
	return nil
}

// 更新分类文章计数
func updateKindEssayCount(tx *sqlx.Tx, kid int, ifAdd bool) error {
	var sqlStr string
	sqlStr = `
        UPDATE kind 
        SET essay_count = essay_count  - 1
        WHERE id = ?`

	if ifAdd {
		sqlStr = `
        UPDATE kind 
        SET essay_count = essay_count + 1
        WHERE id = ?`
	}

	_, err := tx.Exec(sqlStr, kid)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEssay(id int) error {
	return withTx(func(tx *sqlx.Tx) error {
		if err := deleteLabels(tx, id); err != nil {
			return fmt.Errorf("deleteLabels failed,err:%w", err)
		}
		kid, err := deleteEssay(tx, id)
		if err != nil {
			return fmt.Errorf("deleteEssay failed,err:%w", err)
		}

		if err := updateKindEssayCount(tx, kid, false); err != nil {
			return fmt.Errorf("updateKindEssayCount failed,err:%w", err)
		}

		return nil
	})
}

func deleteEssay(tx *sqlx.Tx, eid int) (kid int, err error) {
	sqlStr1 := `SELECT kind_id FROM essay WHERE id = ?`
	sqlStr2 := `DELETE FROM essay WHERE id = ?`

	if err = tx.QueryRow(sqlStr1, eid).Scan(&kid); err != nil {
		return 0, err
	}

	result, err := tx.Exec(sqlStr2, eid)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err = noAffectedRowErr(affected, err, essayNoExist); err != nil {
		return 0, err
	}
	return kid, err
}

func deleteLabels(tx *sqlx.Tx, eid int) error {
	sqlStr := `DELETE FROM essay_label WHERE essay_id = ?`
	if _, err := tx.Exec(sqlStr, eid); err != nil {
		return err
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
