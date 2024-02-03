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
	sqlStr := `SELECT COUNT(*) FROM classifyKind WHERE ClassifyKindName = ?`
	err = db.Get(&row, sqlStr, c.ClassifyKind)
	if row == 0 {
		//创建这个classifyKind
		sqlStr = `INSERT INTO classifyKind(ClassifyKindName) VALUES (?)`
		_, err = db.Exec(sqlStr, c.ClassifyKind)
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
	sqlStr := `SELECT COUNT(*) FROM classify WHERE ClassifyName = ?`
	var count int
	err = db.Get(&count, sqlStr, c.ClassifyName)
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
	//1.先查ClassifyKind里面有没有这个分类 没有就创建 有的话就返回
	err = CheckClassifyKindExist(c)
	if err != nil {
		return err
	}
	//2.再查classify里面有没有这个classifyName
	err = CheckClassifyExist(c)
	if err != nil {
		return err
	}
	//3.在classify表里面添加数据
	sqlStr := `INSERT INTO classify(ClassifyKindName, ClassifyName,ClassifyRoute) VALUES(?,?,?)`
	_, err = db.Exec(sqlStr, c.ClassifyKind, c.ClassifyName, c.ClassifyRoute)
	if err != nil {
		return err
	}
	return
}

// CheckEssayExist 检测文章名称是否存在
func CheckEssayExist(c *models.EssayParams) (err error) {
	sqlStr := `SELECT COUNT(*) FROM essay WHERE  essayKind = ? AND  essayName = ? `
	var count int
	err = db.Get(&count, sqlStr, c.EssayKind, c.EssayName)
	fmt.Println(count)
	if count > 0 {
		return errors.New(essayExist)
	}
	return err
}

// CreateEssay 添加新文章
func CreateEssay(e *models.EssayParams) (err error) {
	//1.检测该文章是否已经存在
	err = CheckEssayExist(e)
	if err != nil {
		return err
	}
	//2.添加该文章
	sqlStr := `INSERT INTO essay(essayKind,essayName, essayContent,essayRoute) values(?,?,?,?)`
	_, err = db.Exec(sqlStr, e.EssayKind, e.EssayName, e.EssayContent, e.EssayRoute)
	return err
}
