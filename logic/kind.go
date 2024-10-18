package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

// AddClassify 新增分类逻辑
func AddClassify(c *models.ClassifyParams) (err error) {
	//1.先查ClassifyKind里面有没有这个分类 没有就创建 有的话就返回
	if err = mysql.CheckClassifyKindExist(c); err != nil {
		return err
	}
	//2.再查classify里面有没有这个classifyName
	if err = mysql.CheckClassifyExist(c); err != nil {
		return err
	}
	//3.在classify表里面添加数据
	return mysql.AddClassify(c)
}

// UpdateKind 更新总纲逻辑
func UpdateKind(k *models.UpdateKindParams) (err error) {
	//1.根据id查询到oldKind
	var oldName string
	if oldName, err = mysql.CheckKind(k.ID); err != nil {
		return err
	}
	//2.更新Kind的kind和classify的name
	return mysql.UpdateKind(oldName, k)
}
