package help

import "blog/dao/redis"

func GetMaliciousIp(maliciousIpMap map[string]bool) error {
	pips, err := redis.GetAllMaliciousIp()
	if err != nil {
		return err
	}
	for _, ip := range *pips {
		maliciousIpMap[ip] = true
	}
	return err
}
