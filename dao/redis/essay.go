package redis

import (
	"blog/dao/mysql"
	"fmt"
)

const (
	scoreIncrement = 1
)

func DeleteEssay(id int) (err error) {
	changeKey := getRedisKey(KeyChangeVisitedTimes) //SMembers
	visitedKey := getRedisKey(KeyVisitedTimes)
	keywordKey := getRedisKey(KeyEssayKeyword)
	eid, err := mysql.GetEssaySnowflakeID(id)
	if err != nil {
		return err
	}
	//1. 删除访问次数
	if err = client.SRem(changeKey, eid).Err(); err != nil {
		return err
	}
	if err := client.HDel(visitedKey, fmt.Sprintf("%d", eid)).Err(); err != nil {
		return err
	}
	//2. 删除关键字
	if err := client.Del(fmt.Sprintf("%s%d", keywordKey, eid), "*").Err(); err != nil {
		return err
	}
	return
}
