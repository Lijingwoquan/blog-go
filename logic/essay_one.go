package logic

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
	"fmt"
)

// GetEssayData 得到文章数据
func GetEssayData(data *models.EssayContent, id int) error {
	return mysql.GetEssayData(data, id)
}

// CreateEssay 新增文章逻辑
func CreateEssay(e *models.EssayParams) error {
	return mysql.CreateEssay(e)
}

// DeleteEssay 删除文章逻辑
func DeleteEssay(id int) error {
	//从mysql里面删除该文章
	if err := redis.DeleteEssay(id); err != nil {
		return err
	}

	//删除redis中文章的相关数据
	return mysql.DeleteEssay(id)
}

// UpdateEssayMsg 更新文章逻辑
func UpdateEssayMsg(u *models.UpdateEssayMsgParams) error {
	//更新数据
	if err := mysql.UpdateEssayMsg(u); err != nil {
		fmt.Println(err)
		return err
	}
	iDAndKeywords := new(models.EssayIdAndKeyword)
	iDAndKeywords.EssayId = u.Id
	iDAndKeywords.Keywords = u.Keywords
	return redis.SetEssayKeyword(iDAndKeywords)
}
