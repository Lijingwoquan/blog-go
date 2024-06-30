package logic

import (
	"blog/dao/mysql"
	"blog/dao/redis"
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

// CreateEssay 新增文章逻辑
func CreateEssay(e *models.EssayParams) (err error) {
	//1.检测该文章是否已经存在
	if err = mysql.CheckEssayExist(e); err != nil {
		return err
	}
	//2.添加该文章
	var eid int64
	//mysql处理数据
	if eid, err = mysql.CreateEssay(e); err != nil {
		return err
	}
	//redis处理数据 --> 初始访问次数
	if err = redis.InitVisitedTimes(eid); err != nil {
		return err
	}
	return
}

// UpdateEssayMsg 更新文章逻辑
func UpdateEssayMsg(u *models.UpdateEssayMsg) error {
	//更新数据
	return mysql.UpdateEssayMsg(u)
}

func UpdateEssayContent(u *models.UpdateEssayContent) error {
	//更新数据
	return mysql.UpdateEssayContent(u)
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

// DeleteEssay 删除文章逻辑
func DeleteEssay(id int) error {
	//从mysql里面删除该文章
	return mysql.DeleteEssay(id)
}
