package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

// UpdateClassify 更新分类信息逻辑
func UpdateClassify(c *models.UpdateClassifyParams) (err error) {
	//1.先由id查询得到oldName
	var oldName string
	if oldName, err = mysql.CheckClassifyName(c.ID); err != nil {
		return err
	}
	//2.由传进来的id查询数据库 进行更新
	return mysql.UpdateClassify(oldName, c)
}
