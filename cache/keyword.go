package cache

import (
	"blog/dao/redis"
)

var (
	GlobalEssayKeyword    []string
	GlobalEssayKeywordMap = make(map[string][]string)
)

func initEssayKeyword() (err error) {
	GlobalEssayKeyword, err = redis.InitEssayKeyWord(GlobalEssayKeywordMap)
	return err
}
