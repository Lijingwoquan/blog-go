package redis

import (
	"blog/dao/mysql"
	"blog/models"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

const (
	year      = "year"
	month     = "month"
	week      = "week"
	yearTime  = time.Hour * 24 * 7 * 12
	monthTime = time.Hour * 24 * 7 * 30
	weekTime  = time.Hour * 24 * 7
)

func SetEssayKeyword(essayKeyword *models.EssayIdAndKeyword) (err error) {
	var eid int64
	if eid, err = mysql.GetEssaySnowflakeID((*essayKeyword).EssayId); err != nil {
		return err
	}

	key := fmt.Sprintf("%s%d", getRedisKey(KeyEssayKeyword), eid)

	// 创建 Redis 管道
	pipe := client.Pipeline()

	// 首先删除现有的所有关键词
	pipe.Del(key)

	// 如果有新的关键词，则设置它们
	if len((*essayKeyword).Keywords) > 0 {
		// 使用 SADD 命令添加到集合
		for _, keyword := range (*essayKeyword).Keywords {
			pipe.SAdd(key, strings.ToLower(strings.TrimSpace(keyword)))
		}
	}

	// 执行管道命令
	_, err = pipe.Exec()
	if err != nil {
		return fmt.Errorf("failed to set essay keywords: %w", err)
	}
	return nil
}

func IncreaseSearchKeyword(keyword string) (err error) {
	return SetYearMonthWeekTimesZoneForZset(keyword, KeyEssayKeyword, 1)
}

func GetEssayKeywordsForIndex(e *[]models.DataAboutEssay) (err error) {
	keyPre := getRedisKey(KeyEssayKeyword)
	for i := range *e {
		ids, err := mysql.GetEssaySnowflakeID((*e)[i].ID)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s%d", keyPre, ids)
		keywords, err := client.SMembers(key).Result()
		if err != nil {
			return err
		}
		(*e)[i].Keywords = keywords
	}
	return err
}

func GetEssayKeywordsForOne(e *models.EssayData) (err error) {
	keyPre := getRedisKey(KeyEssayKeyword)
	key := fmt.Sprintf("%s%d", keyPre, e.Eid)
	keywords, err := client.SMembers(key).Result()
	if err != nil {
		return err
	}
	(*e).Keywords = keywords
	return err
}

func GetSearchKeywordRank(rankKind *models.KeywordRankKind) (err error) {
	//	得到年月日的keywords的zset
	yearKey := fmt.Sprintf("%s%s:", getRedisKey(KeySearchKeyWordTimes), year)
	monthKey := fmt.Sprintf("%s%s:", getRedisKey(KeySearchKeyWordTimes), month)
	weekKey := fmt.Sprintf("%s%s:", getRedisKey(KeySearchKeyWordTimes), week)

	// 从每个zset中获取前10条数据
	yearList, err := getTop10FromZSet(client, yearKey)
	if err != nil {
		return err
	}

	monthList, err := getTop10FromZSet(client, monthKey)
	if err != nil {
		return err
	}

	weekList, err := getTop10FromZSet(client, weekKey)
	if err != nil {
		return err
	}

	// 合并结果
	*rankKind = models.KeywordRankKind{
		Year:  yearList,
		Month: monthList,
		Week:  weekList,
	}

	return nil
}

func getTop10FromZSet(rdb *redis.Client, key string) ([]models.KeywordRankList, error) {
	result, err := rdb.ZRevRangeWithScores(key, 0, 9).Result()
	if err != nil {
		return nil, err
	}

	var rankList []models.KeywordRankList
	for _, z := range result {
		rankList = append(rankList, models.KeywordRankList{
			Keyword: z.Member.(string),
			Times:   int(z.Score),
		})
	}

	return rankList, nil
}
