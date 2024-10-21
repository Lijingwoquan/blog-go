package mysql

import (
	"blog/models"
	"blog/utils"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"sync"
)

const (
	invalidLabelIds = "labels参数无效"
	essayNoExist    = "该文章不存在"
)

func GetEssayData(data *models.EssayContent, id int) error {
	if err := getEssay(data, id); err != nil {
		return fmt.Errorf("getEssay failed,err:%w", err)
	}

	data.NearEssayList = make([]models.EssayData, 0, 5)
	if err := getNearbyEssays(&data.NearEssayList, data.KindID, id); err != nil {
		return fmt.Errorf("getNearbyEssays failed,err:%w", err)
	}
	return nil
}

func getEssay(data *models.EssayContent, id int) error {
	var wg sync.WaitGroup
	wg.Add(2)
	var errChan = make(chan error, 2)
	go func() {
		defer wg.Done()
		if err := getEssayContent(data, id); err != nil {
			errChan <- fmt.Errorf("getEssayContent failed,err:%w", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := increaseEssayCount(id); err != nil {
			errChan <- fmt.Errorf("increaseEssayCount failed,err:%w", err)
			return
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

func getEssayContent(data *models.EssayContent, id int) error {
	sqlStr := `
		SELECT e.id,e.name,e.kind_id, e.content, e.introduction, e.created_time, e.visited_times,
			k.name AS kind_name
		FROM essay e
		LEFT JOIN kind k on e.kind_id = k.id
		where e.id = ?
		`
	return db.Get(data, sqlStr, id)
}

func increaseEssayCount(id int) error {
	sqlStr := `
	UPDATE essay SET visited_times = visited_times + 1
		WHERE id = ?`
	_, err := db.Exec(sqlStr, id)
	return err
}

func getNearbyEssays(data *[]models.EssayData, kID int, eID int) error {
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
	// 构建批量插入的SQL和参数
	sqlStr := `
        INSERT INTO essay_label (essay_id, label_id, label_name)
        SELECT ?, l.id, l.name 
        FROM label l
        WHERE l.id IN (?)
    `

	// 通过sqlx.In来处理IN查询
	query, args, err := sqlx.In(sqlStr, eid, lIDs)
	if err != nil {
		return err
	}

	// 将SQL转换为底层数据库驱动可执行的格式
	query = tx.Rebind(query)

	result, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}

	// 检查影响的行数是否符合预期
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// 如果影响的行数与传入的标签ID数量不匹配,说明有无效的标签ID
	if int(affected) != len(lIDs) {
		return fmt.Errorf(invalidLabelIds)
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

func UpdateEssay(data *models.EssayParams) error {
	return withTx(func(tx *sqlx.Tx) error {
		var err error
		var okid int
		// 先查essay表 得到kind_id
		if okid, err = getEssayKindID(tx, data.ID); err != nil {
			return fmt.Errorf("getEssayKindID failed,err:%w", err)
		}
		// 判断是否需要更新kind表
		if okid != data.KindID {
			// 原来的kind减1
			if err = updateKindEssayCount(tx, okid, false); err != nil {
				return fmt.Errorf("updateKindEssayCount failed,err:%w", err)
			}
			// 新的kind加1
			if err = updateKindEssayCount(tx, data.KindID, true); err != nil {
				return fmt.Errorf("updateKindEssayCount failed,err:%w", err)
			}
		}

		// 更新essay_label
		// 删除原来的essay_label表关联的数据 然后重建?

		// 更新essay表
		if err = updateEssay(tx, data); err != nil {
			return fmt.Errorf("updateEssay failed,err%w", err)
		}
		return err
	})

}

func getEssayKindID(tx *sqlx.Tx, eid int) (kid int, err error) {
	sqlStr := `SELECT kind_id FROM essay WHERE id = ?`

	if err := tx.QueryRow(sqlStr, eid).Scan(&kid); err != nil {
		return 0, err
	}

	return kid, nil
}

func updateEssay(tx *sqlx.Tx, data *models.EssayParams) error {
	sqlStr := `UPDATE essay SET 
               name = :name,
               kind_id = :kind_id,
               introduction = :introduction,
               content = :content,
               img_url = :img_url,
               if_top = :if_top,
               if_recommend = :if_recommend
               WHERE id = :id`
	_, err := tx.NamedExec(sqlStr, data)
	return err
}
