package redis

import (
	"blog/dao/mysql"
	"fmt"
	"strconv"
)

func GetVisitedTimes(eid string) (int64, error) {
	var vt int64
	key := getRedisKey(KeyVisitedTimes)
	//先检查键是否存在
	exist, err := client.HExists(key, eid).Result()
	if err != nil {
		return 0, fmt.Errorf("exist,err := client.HExists(pre,eid).Result() failed,err:%v", err)
	}
	if !exist {
		vt, err = mysql.GetVisitedTimesFromMySQL(eid)
		if err != nil {
			return 0, fmt.Errorf("err,vt = mysql.GetHashValue() failed,err:%v", err)
		}
	}

	vt, err = client.HIncrBy(key, eid, 1).Result()
	if err != nil {
		return 0, fmt.Errorf("get visited times failed: %v", err)
	}

	//将访问的文章加入到集合中
	_, err = client.SAdd(getRedisKey(KeyChangeVisitedTimes), eid).Result()
	if err != nil {
		return 0, fmt.Errorf("client.SAdd(getRedisKey(KeyChangeVisitedTimes),eid).Result() failed,err:%v", err)
	}
	return vt, nil
}

func InitVisitedTimes(eid int64) error {
	pre := getRedisKey(KeyVisitedTimes)
	_, err := client.HSet(pre, fmt.Sprintf("%d", eid), 0).Result()
	return err
}

func GetChangedVisitedTimes() (map[int64]int64, error) {
	pre := getRedisKey(KeyVisitedTimes)
	eids, err := client.SMembers(getRedisKey(KeyChangeVisitedTimes)).Result()
	if err != nil {
		return nil, fmt.Errorf("client.SMembersMap(getRedisKey(KeyChangeVisitedTimes)).Result() failed,err:%v", err)
	}
	var visitedTimesChangedMap = make(map[int64]int64)
	for _, eid := range eids {
		eida, err := strconv.ParseInt(eid, 10, 64)
		vt, err := client.HGet(pre, eid).Int64()
		if err != nil {
			return nil, fmt.Errorf("vt,err := client.HGet(pre,eid).Int64() failed,err:%v", err)
		}
		visitedTimesChangedMap[eida] = vt
	}
	//删除这个集合
	_, err = client.Del(getRedisKey("visited:eids")).Result()
	if err != nil {
		return nil, fmt.Errorf("client.Del(getRedisKey(\"visited:eids\")).Result() failed,err:%v", err)
	}
	return visitedTimesChangedMap, nil
}
