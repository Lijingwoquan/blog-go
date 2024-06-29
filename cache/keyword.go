package cache

import (
	"blog/dao/redis"
)

var (
	GlobalEssayKeyword    []string
	GlobalEssayKeywordMap map[string][]string
)

func initEssayKeyword() (err error) {
	GlobalEssayKeyword, err = redis.InitEssayKeyWord(GlobalEssayKeywordMap)
	return err
}
