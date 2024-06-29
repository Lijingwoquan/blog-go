package redis

import (
	"blog/dao/mysql"
	"blog/models"
	"fmt"
	"strings"
	"time"
)

const (
	year           = "year"
	month          = "month"
	week           = "week"
	yearTime       = time.Hour * 24 * 7 * 12
	monthTime      = time.Hour * 24 * 7 * 30
	weekTime       = time.Hour * 24 * 7
	ScoreIncrement = 1
)

func InitEssayKeyWord(GlobalEssayKeywordMap map[string][]string) (keySlice []string, err error) {
	pattern := "blog:essay:keyword*"
	keys, err := client.Keys(pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("client.Keys(pattern).Result() failed:%w", err)
	}

	for _, key := range keys {
		// 使用 HGETALL 获取键的所有字段和值
		values, err := client.HGetAll(key).Result()
		if err != nil {
			return nil, fmt.Errorf("Error getting values for key %s: %w\n", key, err)
		}
		//blog:essay:keyword:40784933005824000
		//eids := strings.Split(key, ":")[3]
		//eid,_ := strconv.ParseInt(eids,10,64)

		// 处理键的内容
		for field, value := range values {
			//GlobalEssayKeywordMap[eids] = [...GlobalEssayKeywordMap[1],value]
			fmt.Printf("failed :%v \t value:%v", field, value)
		}
	}
	return keySlice, err
}

func SetEssayKeyword(essayKeyword *models.EssayKeyword) (err error) {
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
	// 1.首先对关键词做去空格和小写化
	initKeyword := strings.ToLower(strings.TrimSpace(keyword))

	// 2.给每个关键词一个总的统计次数 分别实现 年 月 周 关键词统计
	yearKey := fmt.Sprintf("%s%s:", getRedisKey(KeySearchKeyWordTimes), year)
	monthKey := fmt.Sprintf("%s%s:", getRedisKey(KeySearchKeyWordTimes), month)
	weekKey := fmt.Sprintf("%s%s:", getRedisKey(KeySearchKeyWordTimes), week)

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
