package redis

import (
	"blog/models"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

const (
	rankCount int64 = 5
)

type RemainingTime struct {
	Year  int64 // 秒
	Month int64 // 秒
	Week  int64 // 秒
}

func getRemainingTime() RemainingTime {
	shanghaiLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		shanghaiLocation = time.Local
	}

	now := time.Now().In(shanghaiLocation)
	currentYear, currentMonth, _ := now.Date()

	// 计算年末时间
	endOfYear := time.Date(currentYear+1, time.January, 1, 0, 0, 0, 0, shanghaiLocation)

	// 计算下个月的第一天
	firstOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, shanghaiLocation)

	// 计算本周末时间
	daysUntilEndOfWeek := time.Saturday - now.Weekday()
	if daysUntilEndOfWeek <= 0 {
		daysUntilEndOfWeek += 7
	}
	endOfWeek := now.Add(time.Duration(daysUntilEndOfWeek*24) * time.Hour).Truncate(24 * time.Hour)

	return RemainingTime{
		Year:  int64(endOfYear.Sub(now).Seconds()),
		Month: int64(firstOfNextMonth.Sub(now).Seconds()),
		Week:  int64(endOfWeek.Sub(now).Seconds()),
	}
}

// SetYearMonthWeekTimesZoneForZset 设置年月日相关Zset
func SetYearMonthWeekTimesZoneForZset(preKey string, param string, scoreIncrement float64) (err error) {
	// 1.首先对成员做去空格和小写化
	member := strings.ToLower(strings.TrimSpace(param))

	// 2.给每个成员一个总的统计次数 分别实现 年 月 周 关键词统计
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 3.得到剩余时间
	remainingTime := getRemainingTime()

	// 4.用集合实现 --> 内置排序
	pipe := client.Pipeline()

	// 年统计
	pipe.ZIncrBy(yearKey, scoreIncrement, member)
	pipe.Expire(yearKey, time.Duration(remainingTime.Year)*time.Second)

	// 月统计
	pipe.ZIncrBy(monthKey, scoreIncrement, member)
	pipe.Expire(monthKey, time.Duration(remainingTime.Month)*time.Second)

	// 周统计
	pipe.ZIncrBy(weekKey, scoreIncrement, member)
	pipe.Expire(weekKey, time.Duration(remainingTime.Week)*time.Second)

	// 执行管道命令
	if _, err = pipe.Exec(); err != nil {
		return fmt.Errorf("failed to increase Zset score: %w", err)
	}
	return nil
}

// SetYearMonthWeekTimesZoneForSet 设置年月日相关set
func SetYearMonthWeekTimesZoneForSet(preKey string, param string) (err error) {
	// 1.首先对关键词做去空格和小写化
	member := strings.ToLower(strings.TrimSpace(param))

	// 2.给每个关键词一个总的统计次数 分别实现 年 月 周 关键词统计
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 3.得到剩余时间
	remainingTime := getRemainingTime()

	// 4.用集合实现 --> 内置排序
	pipe := client.Pipeline()

	// 年统计
	pipe.SAdd(yearKey, member)
	pipe.Expire(yearKey, time.Duration(remainingTime.Year)*time.Second)

	// 月统计
	pipe.SAdd(monthKey, member)
	pipe.Expire(monthKey, time.Duration(remainingTime.Month)*time.Second)

	// 周统计
	pipe.SAdd(weekKey, member)
	pipe.Expire(weekKey, time.Duration(remainingTime.Week)*time.Second)

	// 执行管道命令
	if _, err = pipe.Exec(); err != nil {
		return fmt.Errorf("failed to set member: %w", err)
	}
	return nil
}

// GetYearMonthWeekTimesZoneForZsetRank 得到年月日相关Zset
func GetYearMonthWeekTimesZoneForZsetRank(rankKind *models.RankKindForZset, preKey string) (err error) {
	//	得到年月日的keywords的zset
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 从每个zset中获取前10条数据
	yearList, err := getTopXFromZSet(yearKey, rankCount)
	if err != nil {
		return err
	}

	monthList, err := getTopXFromZSet(monthKey, rankCount)
	if err != nil {
		return err
	}

	weekList, err := getTopXFromZSet(weekKey, rankCount)
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

// getTopXFromZSet 得到Zset排序
func getTopXFromZSet(key string, count int64) (models.RankListForZset, error) {
	result, err := client.ZRevRangeWithScores(key, 0, count).Result()
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

// GetYearMonthWeekTimesZoneForSet 得到年月日相关set
func GetYearMonthWeekTimesZoneForSet(setKind *models.UserIpForSet, preKey string) (err error) {
	//	得到年月日的keywords的zset
	yearKey := fmt.Sprintf("%s%s:", preKey, year)
	monthKey := fmt.Sprintf("%s%s:", preKey, month)
	weekKey := fmt.Sprintf("%s%s:", preKey, week)

	// 从每个zset中获取前10条数据
	yearCount, err := getCountFromSet(yearKey)
	if err != nil {
		return err
	}

	monthCount, err := getCountFromSet(monthKey)
	if err != nil {
		return err
	}

	weekCount, err := getCountFromSet(weekKey)
	if err != nil {
		return err
	}

	// 合并结果
	*setKind = models.UserIpForSet{
		Year:  yearCount,
		Month: monthCount,
		Week:  weekCount,
	}

	return nil
}

// getCountFromSet 得到set的成员值
func getCountFromSet(key string) (int64, error) {
	count, err := client.SCard(key).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ticker

// CleanLowerZsetEveryMonth 删除低频元素
func CleanLowerZsetEveryMonth() error {
	preKey := getRedisKey(KeySearchKeyWordTimes)
	yearKey := fmt.Sprintf("%s%s:", preKey, year)

	// 获取分数小于等于5的成员
	lowFrequentMembers, err := client.ZRangeByScoreWithScores(yearKey, redis.ZRangeBy{
		Min: "-inf",
		Max: "5",
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to get members with score <= 5 from %s: %w", yearKey, err)
	}

	if len(lowFrequentMembers) > 0 {
		// 提取成员名称
		var membersToRemove []interface{}
		for _, z := range lowFrequentMembers {
			membersToRemove = append(membersToRemove, z.Member)
		}

		// 删除这些成员
		if err = client.ZRem(yearKey, membersToRemove...).Err(); err != nil {
			return fmt.Errorf("failed to remove low frequent members from %s: %w", yearKey, err)
		}
	}

	return nil
}
