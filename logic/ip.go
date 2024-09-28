package logic

import (
	"blog/cache"
	"blog/dao/redis"
	"fmt"
	"go.uber.org/zap"
)

const (
	preSecondRequest = 15
	maliciousTimes   = 100
)

const (
	requestTooFrequent = "请求过于频繁,请稍后再试"
	ipForbid           = "恭喜你!ip已被永久封禁"
)

func SaveUserIp(ip string) error {
	if err := redis.SaveUserIp(ip); err != nil {
		return err
	}
	return nil
}

func IpLimit(ip string) (err error) {
	// 1.检查ip是否在恶意请求里
	var exist bool
	exist = cache.IfIpInMaliciousMap(ip)
	if exist {
		return fmt.Errorf(ipForbid)
	}

	// 2.增加ip单位时间请求次数
	var times int64
	if times, err = redis.IncreaseIpRequestTimes(ip); err != nil {
		zap.L().Error("redis.IncreaseRequestTimes(ip),err:%v", zap.Error(err))
		return err
	}

	// 3.是否超出恶意指标处理
	if times > maliciousTimes {
		//在redis里存放
		if err = redis.SetIpMalicious(ip); err != nil {
			zap.L().Error("redis.SetIpMalicious(ip) failed,err:%v", zap.Error(err))
			return err
		}
		//在缓存里存放
		cache.AddMaliciousIpMap(ip)
		return fmt.Errorf(ipForbid)
	}

	// 4.是否频繁请求
	if times > preSecondRequest {
		return fmt.Errorf(requestTooFrequent)
	}

	return err
}
