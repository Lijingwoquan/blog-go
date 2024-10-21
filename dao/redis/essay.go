package redis

import (
	"fmt"
)

func DeleteEssay(id int) (err error) {
	keywordKey := getRedisKey(KeyEssayKeyword)
	if err := client.Del(fmt.Sprintf("%s%d", keywordKey, id), "*").Err(); err != nil {
		return err
	}
	return
}
