package redis

import (
	"blog/models"
	"fmt"
	"strings"
)

const (
	year                 = "year"
	month                = "month"
	week                 = "week"
	defaultIncreaseCount = 1
)

func SetEssayKeyword(essayKeyword *models.EssayIdAndKeyword) (err error) {

	key := fmt.Sprintf("%s%d", getRedisKey(KeyEssayKeyword), essayKeyword.EssayId)

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

func IncreaseSearchKeyword(preKey string, keyword string) (err error) {
	preKey = getRedisKey(preKey)
	return SetYearMonthWeekTimesZoneForZset(preKey, keyword, defaultIncreaseCount)
}

// GetEssayKeywords 获取文章关键字
func GetEssayKeywords(e *[]models.EssayData) (err error) {
	keyPre := getRedisKey(KeyEssayKeyword)
	for i := range *e {
		key := fmt.Sprintf("%s%d", keyPre, (*e)[i].ID)
		keywords, err := client.SMembers(key).Result()
		if err != nil {
			return err
		}
		(*e)[i].Keywords = append(keywords, (*e)[i].Name)
	}
	return err
}

func GetSearchKeywordRank(rankKind *models.RankKindForZset) (err error) {
	preKey := getRedisKey(KeySearchKeyWordTimes)
	return GetYearMonthWeekTimesZoneForZsetRank(rankKind, preKey)
}

func CleanLowerKeywordsZsetEveryMonth() error {
	return CleanLowerZsetEveryMonth()
}
