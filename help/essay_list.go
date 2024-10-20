package help

import (
	"blog/dao/mysql"
	"blog/dao/redis"
	"blog/models"
)

func ResponseDataAboutEssayList(essayList *[]models.EssayData) (err error) {
	if err = mysql.GetAllEssay(essayList); err != nil {
		return err
	}

	return redis.GetEssayKeywords(essayList)
}
