package mysql

import (
	"blog/models"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

const (
	classifyExist = "该分类已存在"
	essayExist    = "该文章已存在"
)

// CheckClassifyKindExist 检查分类种类存不存在
func CheckClassifyKindExist(c *models.ClassifyParams) (err error) {
	var row int
	sqlStr := `SELECT COUNT(*) FROM classifyKind WHERE name = ?`
	err = db.Get(&row, sqlStr, c.Kind)
	if row == 0 {
		//创建这个classifyKind
		sqlStr = `INSERT INTO classifyKind(name) VALUES (?)`
		_, err = db.Exec(sqlStr, c.Kind)
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
	sqlStr := `INSERT INTO essay(kind,name, content,router) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, e.Kind, e.Name, e.Content, e.Router)
	return err
}

// UpdateEssay 更新文章
func UpdateEssay(data *models.UpdateEssay) (err error) {
	sqlStr := `UPDATE essay SET name= ?,kind = ? ,content = ?,router = ? WHERE name = ?`
	result, err := db.Exec(sqlStr, data.Name, data.Kind, data.Content, data.Router, data.OldName)
	if err != nil {
		return
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(essayNotExist)
	}
	return
}
