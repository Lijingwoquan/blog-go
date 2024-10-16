package logic

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
	"fmt"
)

// GetEssayData 得到文章数据
func GetEssayData(data *models.EssayData, id int) error {
	var err error
	if err = mysql.GetEssayData(data, id); err != nil {
		return err
	}
	return nil
}

// CreateEssay 新增文章逻辑
func CreateEssay(e *models.EssayParams) (err error) {
	//1.检测该文章是否已经存在
	if err = mysql.CheckEssayExist(e); err != nil {
		return err
	}

	//mysql处理数据
	if _, err = mysql.CreateEssay(e); err != nil {
		return err
	}

	return
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
