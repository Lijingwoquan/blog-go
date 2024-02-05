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

// UpdateEssay 更新文章逻辑
func UpdateEssay(u *models.UpdateEssay) (err error) {
	//更新数据
	err = mysql.UpdateEssay(u)
	return
}
