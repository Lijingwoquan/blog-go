package mysql

import (
	"blog/models"
	"blog/pkg/snowflake"
	"errors"
	"time"
)

// CheckClassifyKindExist 检查分类种类存不存在
func CheckClassifyKindExist(c *models.ClassifyParams) error {
	var row int
	var err error
	sqlStr := `SELECT COUNT(*) FROM kind WHERE name = ?`
	if err = db.Get(&row, sqlStr, c.Kind); err != nil {
		return err
	}
	if row == 0 {
		//这里就必须携带icon
		if c.Icon == "" {
			return errors.New(needIcon)
		}
		//创建这个classifyKind
		sqlStr = `INSERT INTO kind(name,icon) VALUE (?,?)`
		if _, err = db.Exec(sqlStr, c.Kind, c.Icon); err != nil {
			return err
		}
	}
	return nil
}

// CheckClassifyExist 检查分类存不存在
func CheckClassifyExist(c *models.ClassifyParams) error {
	var err error
	sqlStr1 := `SELECT COUNT(*) FROM classify WHERE name = ?`
	var count int
	if err = db.Get(&count, sqlStr1, c.Name); err != nil {
		return err
	}
	if count > 0 {
		return errors.New(classifyExist)
	}
	count = 0
	sqlStr2 := `SELECT COUNT(*) FROM classify WHERE router = ?`
	if err = db.Get(&count, sqlStr2, c.Router); err != nil {
		return err
	}
	if count > 0 {
		return errors.New(classifyExist)
	}
	return nil
}

// AddClassify 添加新分类
func AddClassify(c *models.ClassifyParams) error {
	var err error
	sqlStr := `INSERT INTO classify(kind, name,router) VALUES(?,?,?)`
	if _, err = db.Exec(sqlStr, c.Kind, c.Name, c.Router); err != nil {
		return err
	}
	return nil
}

// CheckEssayExist 检测文章是否存在
func CheckEssayExist(c *models.EssayParams) error {
	var err error
	sqlStr1 := `SELECT COUNT(*) FROM essay WHERE  kind = ? AND  name = ? `
	var count int
	if err = db.Get(&count, sqlStr1, c.Kind, c.Name); err != nil {
		return err
	}

	if count > 0 {
		return errors.New(essayExist)
	}

	count = 0
	sqlStr2 := `SELECT COUNT(*) FROM essay WHERE router = ?`
	if err = db.Get(&count, sqlStr2, c.Router); err != nil {
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
	if formattedTime, err = getChineseTime(); err != nil {
		return 0, err
	}
	eid := snowflake.GenID()
	sqlStr := `INSERT INTO essay(kind,name,content,router,Introduction,createdTime,updatedTime,eid) values(?,?,?,?,?,?,?,?)`
	if _, err = db.Exec(sqlStr, e.Kind, e.Name, e.Content, e.Router, e.Introduction, formattedTime, formattedTime, eid); err != nil {
		return 0, err
	}
	return eid, err
}

func getChineseTime() (string, error) {
	//加载中国时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	T := time.Now().In(loc)
	t := T.Format("2006-01-02 15:04:05")
	return t, nil
}

// UpdateEssayMsg 更新文章基本信息
func UpdateEssayMsg(data *models.UpdateEssayMsg) error {
	var err error
	var formattedTime string
	if formattedTime, err = getChineseTime(); err != nil {
		return err
	}
	sqlStr := `UPDATE essay SET name= ?,kind = ? ,introduction=?,router = ?,updatedTime=? WHERE id = ?`
	result, err := db.Exec(sqlStr, data.Name, data.Kind, data.Introduction, data.Router, formattedTime, data.Id)
	if err != nil {
		return err
	}
	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return nil
}

// UpdateEssayContent 更新文章内容
func UpdateEssayContent(data *models.UpdateEssayContent) error {
	var err error
	var formattedTime string
	if formattedTime, err = getChineseTime(); err != nil {
		return err
	}
	sqlStr := `UPDATE essay SET content=?,updatedTime=? WHERE id = ?`
	result, err := db.Exec(sqlStr, data.Content, formattedTime, data.Id)
	if err != nil {
		return err
	}
	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return nil
}

func CheckKind(id int) (oldName string, err error) {
	sqlStr := `SELECT name FROM kind WHERE  id = ?`
	err = db.QueryRow(sqlStr, id).Scan(&oldName)
	return
}

func UpdateKind(oldName string, k *models.UpdateKindParams) (err error) {
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
	//更新classify的name
	sqlStr1 := `UPDATE kind SET name =?,icon=? WHERE  id = ?`
	if _, err = tx.Exec(sqlStr1, k.Name, k.Icon, k.ID); err != nil {
		return err
	}
	//更新符合classify的kind
	sqlStr2 := `UPDATE  classify SET kind = ? WHERE kind = ?`
	_, err = tx.Exec(sqlStr2, k.Name, oldName)
	return err
}

func CheckClassifyName(id int) (oldName string, err error) {
	sqlStr := `SELECT name FROM classify WHERE  id = ?`
	err = db.QueryRow(sqlStr, id).Scan(&oldName)
	return
}

func UpdateClassify(oldName string, c *models.UpdateClassifyParams) (err error) {
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
	//更新classify的name
	sqlStr1 := `UPDATE classify SET name=?,router=? WHERE  id = ?`
	_, err = tx.Exec(sqlStr1, c.Name, c.Router, c.ID)
	if err != nil {
		return err
	}
	//更新符合classifyName的essay
	sqlStr2 := `UPDATE  essay SET kind = ? WHERE kind = ?`
	_, err = tx.Exec(sqlStr2, c.Name, oldName)
	return err
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
