package mysql

import (
	"blog/models"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// CheckClassifyKindExist 检查分类种类存不存在
func CheckClassifyKindExist(c *models.ClassifyParams) (err error) {
	var row int
	sqlStr := `SELECT COUNT(*) FROM kind WHERE name = ?`
	err = db.Get(&row, sqlStr, c.Kind)
	if row == 0 {
		//这里就必须携带icon
		if c.Icon == "" {
			return errors.New(needIcon)
		}
		//创建这个classifyKind
		sqlStr = `INSERT INTO kind(name,icon) VALUE (?,?)`
		_, err = db.Exec(sqlStr, c.Kind, c.Icon)
		if err != nil {
			zap.L().Error("db.Exec(sqlStr,c.ClassifyName) failed", zap.Error(err))
			return
		}
	} else {
		return
	}
	return err
}

// CheckClassifyExist 检查分类名存不存在
func CheckClassifyExist(c *models.ClassifyParams) (err error) {
	sqlStr := `SELECT COUNT(*) FROM classify WHERE name = ?`
	var count int
	err = db.Get(&count, sqlStr, c.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(classifyExist)
	}
	return err
}

// AddClassify 添加新分类
func AddClassify(c *models.ClassifyParams) (err error) {
	sqlStr := `INSERT INTO classify(kind, name,router) VALUES(?,?,?)`
	_, err = db.Exec(sqlStr, c.Kind, c.Name, c.Router)
	if err != nil {
		return err
	}
	return
}

// CheckEssayExist 检测文章名称是否存在
func CheckEssayExist(c *models.EssayParams) (err error) {
	sqlStr := `SELECT COUNT(*) FROM essay WHERE  kind = ? AND  name = ? `
	var count int
	err = db.Get(&count, sqlStr, c.Kind, c.Name)
	fmt.Println(count)
	if count > 0 {
		return errors.New(essayExist)
	}
	return err
}

// CreateEssay 添加新文章
func CreateEssay(e *models.EssayParams) (err error) {
	sqlStr := `INSERT INTO essay(kind,name, content,router,Introduction) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, e.Kind, e.Name, e.Content, e.Router, e.Introduction)
	return err
}

// UpdateEssayMsg 更新文章
func UpdateEssayMsg(data *models.UpdateEssayMsg) (err error) {
	updateTime := time.Now()
	formattedTime := updateTime.Format("2006-01-02 15:04:05")
	sqlStr := `UPDATE essay SET name= ?,kind = ? ,router = ?,updatedTime=? WHERE id = ?`
	result, err := db.Exec(sqlStr, data.Name, data.Kind, data.Router, formattedTime, data.Id)
	if err != nil {
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return
}

func UpdateEssayContent(data *models.UpdateEssayContent) (err error) {
	updateTime := time.Now()
	formattedTime := updateTime.Format("2006-01-02 15:04:05")
	sqlStr := `UPDATE essay SET content=?,updatedTime=? WHERE id = ?`
	result, err := db.Exec(sqlStr, data.Content, formattedTime, data.Id)
	if err != nil {
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return
}

func CheckKind(id int) (oldName string, err error) {
	sqlStr := `SELECT name FROM kind WHERE  id = ?`
	err = db.QueryRow(sqlStr, id).Scan(&oldName)
	return
}

func UpdateKind(oldName string, k *models.UpdateKindParams) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
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
	_, err = tx.Exec(sqlStr1, k.Name, k.Icon, k.ID)
	//更新符合classify的kind
	sqlStr2 := `UPDATE  classify SET kind = ? WHERE kind = ?`
	_, err = tx.Exec(sqlStr2, k.Name, oldName)
	return
}

func CheckClassifyName(id int) (oldName string, err error) {
	sqlStr := `SELECT name FROM classify WHERE  id = ?`
	err = db.QueryRow(sqlStr, id).Scan(&oldName)
	return
}

func UpdateClassify(oldName string, c *models.UpdateClassifyParams) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
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
	//更新符合classifyName的essay
	sqlStr2 := `UPDATE  essay SET kind = ? WHERE kind = ?`
	_, err = tx.Exec(sqlStr2, c.Name, oldName)
	return
}

func DeleteEssay(id int) (err error) {
	sqlStr := `DELETE  FROM essay WHERE id=?`
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		return
	}
	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}
	if RowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return
}
