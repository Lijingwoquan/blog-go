package logic

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

func GetEssayData(data *models.EssayContent) error {
	return mysql.GetEssayData(data)
}

// CreateEssay 新增文章逻辑
func CreateEssay(e *models.EssayParams) error {
	return mysql.CreateEssay(e)
	//这里还需要补充redis的添加keyword功能
}

// DeleteEssay 删除文章逻辑
func DeleteEssay(id int) error {
	//删除redis中文章的相关数据
	if err := redis.DeleteEssay(id); err != nil {
		return err
	}

	//从mysql里面删除该文章
	return mysql.DeleteEssay(id)
}

func UpdateEssay(e *models.EssayParams) error {
	//更新数据
	if err := mysql.UpdateEssay(e); err != nil {
		return err
	}
	idKeywords := new(models.EssayIdAndKeyword)
	idKeywords.EssayId = e.ID
	idKeywords.Keywords = e.Keywords
	return redis.SetEssayKeyword(idKeywords)
}
