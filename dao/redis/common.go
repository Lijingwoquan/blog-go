package redis

import (
	"blog/models"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

const (
	rankCount int64 = 10
)

func SetYearMonthWeekTimesZoneForZset(preKey string, param string, scoreIncrement float64) (err error) {
	// 1.首先对关键词做去空格和小写化
	initKeyword := strings.ToLower(strings.TrimSpace(param))

	// 2.给每个关键词一个总的统计次数 分别实现 年 月 周 关键词统计
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 3.用集合实现 --> 内置排序
	pipe := client.Pipeline()

	// 年统计
	pipe.ZIncrBy(yearKey, scoreIncrement, initKeyword)
	pipe.Expire(yearKey, yearTime)

	// 月统计
	pipe.ZIncrBy(monthKey, scoreIncrement, initKeyword)
	pipe.Expire(monthKey, monthTime)

	// 周统计
	pipe.ZIncrBy(weekKey, scoreIncrement, initKeyword)
	pipe.Expire(weekKey, weekTime)

	// 执行管道命令
	if _, err = pipe.Exec(); err != nil {
		return fmt.Errorf("failed to increase search keyword: %w", err)
	}
	return nil
}

func GetYearMonthWeekTimesZoneForZsetRank(rankKind *models.RankKind, preKey string) (err error) {
	//	得到年月日的keywords的zset
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 从每个zset中获取前10条数据
	yearList, err := getTopXFromZSet(client, yearKey, rankCount)
	if err != nil {
		return err
	}

	monthList, err := getTopXFromZSet(client, monthKey, rankCount)
	if err != nil {
		return err
	}

	weekList, err := getTopXFromZSet(client, weekKey, rankCount)
	if err != nil {
		return err
	}

	// 合并结果
	*rankKind = models.RankKind{
		Year:  yearList,
		Month: monthList,
		Week:  weekList,
	}

	return nil
}

func getTopXFromZSet(rdb *redis.Client, key string, count int64) ([]models.RankList, error) {
	result, err := rdb.ZRevRangeWithScores(key, 0, count).Result()
	if err != nil {
		return nil, err
	}

	var rankList []models.RankList
	for _, z := range result {
		rankList = append(rankList, models.RankList{
			Keyword: z.Member.(string),
			Times:   int(z.Score),
		})
	}

	return rankList, nil
}

func SetFrequentZsetItemToSet(preKey string, item string) (err error) {
	preKey = getRedisKey(preKey)
	return client.SAdd(preKey, item).Err()
}

func CleanLowerZsetEveryMonth(key string) error {
	setMembers, err := client.SMembers(key).Result()
	if err != nil {
		return err
	}
	//从年中删除该keyword
	preKey := getRedisKey(KeySearchKeyWordTimes)
	yearKey := fmt.Sprintf("%s%s:", preKey, year)

	zsetMembers, err := client.ZRange(yearKey, 0, -1).Result()
	if err != nil {
		return err
	}
	setMap := make(map[string]struct{}, len(setMembers))
	for _, member := range setMembers {
		setMap[member] = struct{}{}
	}
	var toDelete []string
	for _, member := range zsetMembers {
		if _, exists := setMap[member]; !exists {
			toDelete = append(toDelete, member)
		}
	}

	if len(toDelete) > 0 {
		if err = client.ZRem(yearKey, toDelete).Err(); err != nil {
			return err
		}
	}

	return nil
}
