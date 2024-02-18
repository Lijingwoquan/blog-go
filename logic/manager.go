package logic

import (
	"blog/dao/mysql"
	"blog/models"
)

// AddClassify 新增分类逻辑
func AddClassify(c *models.ClassifyParams) (err error) {
	//1.先查ClassifyKind里面有没有这个分类 没有就创建 有的话就返回
	err = mysql.CheckClassifyKindExist(c)
	if err != nil {
		return err
	}
	//2.再查classify里面有没有这个classifyName
	err = mysql.CheckClassifyExist(c)
	if err != nil {
		return err
	}
	//3.在classify表里面添加数据
	err = mysql.AddClassify(c)
	return
}

// CreateEssay 新增文章逻辑
func CreateEssay(e *models.EssayParams) (err error) {
	//1.检测该文章是否已经存在
	err = mysql.CheckEssayExist(e)
	if err != nil {
		return err
	}
	//2.添加该文章
	err = mysql.CreateEssay(e)
	return
}

// UpdateEssayMsg 更新文章逻辑
func UpdateEssayMsg(u *models.UpdateEssayMSg) (err error) {
	//更新数据
	err = mysql.UpdateEssayMsg(u)
	return
}

// UpdateKind 更新总纲逻辑
func UpdateKind(k *models.UpdateKindParams) (err error) {
	//1.根据id查询到oldKind
	oldName, err := mysql.CheckKind(k.ID)
	if err != nil {
		return err
	}
	//2.更新Kind的kind和classify的name
	err = mysql.UpdateKind(oldName, k)
	return
}

// UpdateClassify 更新分类信息逻辑
func UpdateClassify(c *models.UpdateClassifyParams) (err error) {
	//1.先由id查询得到oldName
	oldName, err := mysql.CheckClassifyName(c.ID)
	if err != nil {
		return
	}
	//2.由传进来的id查询数据库 进行更新
	err = mysql.UpdateClassify(oldName, c)
	return
}

// DeleteEssay 删除文章逻辑
func DeleteEssay(id int) (err error) {
	//1.从mysql里面删除该文章
	err = mysql.DeleteEssay(id)
	return
}
