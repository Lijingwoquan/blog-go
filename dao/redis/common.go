package redis

import (
	"fmt"
	"strings"
)

func SetYearMonthWeekTimesZoneForZset(param string, preKey string, scoreIncrement float64) (err error) {
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
