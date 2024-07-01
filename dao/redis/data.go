package redis

import (
	"blog/dao/mysql"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

const (
	scoreIncrement = 1
)

func GetVisitedTimes(id int) (int64, error) {
	changeKey := getRedisKey(KeyChangeVisitedTimes)
	visitedKey := getRedisKey(KeyVisitedTimes)
	eid, err := mysql.GetEssaySnowflakeID(id)
	eids := fmt.Sprintf("%d", eid)
	if err != nil {
		return 0, err
	}
	var finalVT int64

	txf := func(tx *redis.Tx) error {
		//先检查键是否存在
		exist, err := client.HExists(visitedKey, eids).Result()
		if err != nil {
			return fmt.Errorf("client.HExists(key, eid).Result() failed: %w", err)
		}
		var vt int64
		if !exist {
			vt, err = mysql.GetVisitedTimesFromMySQL(eids)
			if err != nil {
				return fmt.Errorf("mysql.GetVisitedTimesFromMySQL(eid) failed: %w", err)
			}
		}

		if vt, err = tx.HIncrBy(visitedKey, eids, vt+scoreIncrement).Result(); err != nil {
			return fmt.Errorf("tx.HIncrBy(key, eid, scoreIncrement).Result() failed: %w", err)
		}

		// 将访问的文章加入到集合中
		if _, err = tx.SAdd(changeKey, eids).Result(); err != nil {
			return fmt.Errorf("tx.SAdd(key, eid).Result() failed: %w", err)
		}

		finalVT = vt // 保存最终的访问次数
		return nil
	}

	// 执行事务
	if err := client.Watch(txf, visitedKey); err != nil {
		return 0, fmt.Errorf("client.Watch(txf, visitedKey): %w", err)
	}

	return finalVT, nil
}

func InitVisitedTimes(eid int64) error {
	pre := getRedisKey(KeyVisitedTimes)

	_, err := client.HSet(pre, fmt.Sprintf("%d", eid), 0).Result()
	return err
}

func GetAndClearChangedVisitedTimes() (map[int64]int64, error) {
	changeKey := getRedisKey(KeyChangeVisitedTimes)
	visitedKey := getRedisKey(KeyVisitedTimes)

	// 获取发生变化的 eids
	eids, err := client.SMembers(changeKey).Result()
	if err != nil {
		return nil, err
	}

	if len(eids) == 0 {
		return make(map[int64]int64), nil
	}

	visitedTimesChangedMap := make(map[int64]int64, len(eids))
	pipe := client.Pipeline()
	for _, eid := range eids {
		pipe.HGet(visitedKey, eid)
	}

	cmders, err := pipe.Exec()
	if err != nil {
		return nil, fmt.Errorf("pipe.Exec() failed: %w", err)
	}

	for i, cmder := range cmders {
		eidInt, err := strconv.ParseInt(eids[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse eid %s: %w", eids[i], err)
		}
		vt, err := cmder.(*redis.StringCmd).Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to get value for eid %s: %w", eids[i], err)
		}

		visitedTimesChangedMap[eidInt] = vt
	}

	// 删除这个集合
	_, err = client.Del(changeKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to delete key %s: %w", changeKey, err)
	}
	return visitedTimesChangedMap, nil
}

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
