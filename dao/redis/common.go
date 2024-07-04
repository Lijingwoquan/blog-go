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

// Zset

func SetYearMonthWeekTimesZoneForZset(preKey string, param string, scoreIncrement float64) (err error) {
	// 1.首先对关键词做去空格和小写化
	member := strings.ToLower(strings.TrimSpace(param))

	// 2.给每个关键词一个总的统计次数 分别实现 年 月 周 关键词统计
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 3.用集合实现 --> 内置排序
	pipe := client.Pipeline()

	// 年统计
	pipe.ZIncrBy(yearKey, scoreIncrement, member)
	pipe.Expire(yearKey, yearTime)

	// 月统计
	pipe.ZIncrBy(monthKey, scoreIncrement, member)
	pipe.Expire(monthKey, monthTime)

	// 周统计
	pipe.ZIncrBy(weekKey, scoreIncrement, member)
	pipe.Expire(weekKey, weekTime)

	// 执行管道命令
	if _, err = pipe.Exec(); err != nil {
		return fmt.Errorf("failed to increase Zset score: %w", err)
	}
	return nil
}

func GetYearMonthWeekTimesZoneForZsetRank(rankKind *models.RankKindForZset, preKey string) (err error) {
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
	*rankKind = models.RankKindForZset{
		Year:  yearList,
		Month: monthList,
		Week:  weekList,
	}

	return nil
}

func getTopXFromZSet(rdb *redis.Client, key string, count int64) (models.RankListForZset, error) {
	result, err := rdb.ZRevRangeWithScores(key, 0, count).Result()
	if err != nil {
		return models.RankListForZset{}, err
	}

	var rankList models.RankListForZset
	for _, z := range result {
		rankList.X = append(rankList.X, z.Member.(string))
		rankList.Y = append(rankList.Y, int(z.Score))
	}
	return rankList, nil
}

// Set

func SetYearMonthWeekTimesZoneForSet(preKey string, param string) (err error) {
	// 1.首先对关键词做去空格和小写化
	member := strings.ToLower(strings.TrimSpace(param))

	// 2.给每个关键词一个总的统计次数 分别实现 年 月 周 关键词统计
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 3.用集合实现 --> 内置排序
	pipe := client.Pipeline()

	// 年统计
	pipe.SAdd(yearKey, member)
	pipe.Expire(yearKey, yearTime)

	// 月统计
	pipe.SAdd(monthKey, member)
	pipe.Expire(monthKey, monthTime)

	// 周统计
	pipe.SAdd(weekKey, member)
	pipe.Expire(weekKey, weekTime)

	// 执行管道命令
	if _, err = pipe.Exec(); err != nil {
		return fmt.Errorf("failed to set member: %w", err)
	}
	return nil
}

func GetYearMonthWeekTimesZoneForSetCount(setKind *models.SetKind, preKey string) (err error) {
	//	得到年月日的keywords的zset
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 从每个zset中获取前10条数据
	yearCount, err := getCountFromSet(client, yearKey)
	if err != nil {
		return err
	}

	monthCount, err := getCountFromSet(client, monthKey)
	if err != nil {
		return err
	}

	weekCount, err := getCountFromSet(client, weekKey)
	if err != nil {
		return err
	}

	// 合并结果
	*setKind = models.SetKind{
		Year:  yearCount,
		Month: monthCount,
		Week:  weekCount,
	}

	return nil
}

func getCountFromSet(rdb *redis.Client, key string) (int64, error) {
	count, err := rdb.SCard(key).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ticker

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
