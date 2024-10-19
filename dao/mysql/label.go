package mysql

import (
	"blog/models"
	"errors"
)

func GetLabelList(data *[]models.LabelData) error {
	sqlStr := `SELECT name,id FROM label`
	return db.Select(data, sqlStr)
}

func GetOneDataAboutClassify(data *models.LabelData) error {
	sqlStr := `SELECT name,id FROM label WHERE name = ?`
	return db.Get(data, sqlStr, data.Name)
}

// CheckClassifyKindExist 检查分类种类是否存在
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

// CheckClassifyExist 检查分类是否存在
func CheckClassifyExist(c *models.ClassifyParams) error {
	var err error
	sqlStr := `SELECT COUNT(*) FROM classify WHERE name = ?`
	var count int
	if err = db.Get(&count, sqlStr, c.Name); err != nil {
		return err
	}
	if count > 0 {
		return errors.New(classifyExist)
	}
	count = 0
	return nil
}

// AddClassify 添加新分类
func AddClassify(c *models.ClassifyParams) error {
	var err error
	sqlStr := `INSERT INTO classify(kind, name) VALUES(?,?)`
	if _, err = db.Exec(sqlStr, c.Kind, c.Name); err != nil {
		return err
	}
	return nil
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
	sqlStr1 := `UPDATE classify SET name=? WHERE  id = ?`
	_, err = tx.Exec(sqlStr1, c.Name, c.ID)
	if err != nil {
		return err
	}
	//更新符合classifyName的essay
	sqlStr2 := `UPDATE  essay SET kind = ? WHERE kind = ?`
	_, err = tx.Exec(sqlStr2, c.Name, oldName)
	return err
}
