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
	//从mysql查数据
	if err = mysql.GetEssayData(data, id); err != nil {
		return err
	}
	//从redis查数据
	if data.VisitedTimes, err = redis.GetVisitedTimes(fmt.Sprintf("%d", data.Eid)); err != nil {
		return err
	}
	return nil
}
