package logic

import (
	"blog/cache"
	"blog/dao/redis"
	"blog/models"
	"strings"
)

func IncreaseSearchKeyword(SearchKeyword *models.KeywordParam) (err error) {
	preKey := redis.KeySearchKeyWordTimes
	// 向redis中加入keyWord
	return redis.IncreaseSearchKeyword(preKey, (*SearchKeyword).Keyword)
}

func GetEssayListByKeyword(e *[]models.DataAboutEssay, k *models.KeywordParam) {
	for _, essay := range cache.GlobalDataAboutIndex.EssayList {
		// 检查 essay.keyword 数组中是否包含指定的关键字 k
		for _, keyword := range essay.Keywords {
			if strings.Contains(keyword, k.Keyword) {
				// 如果包含，将文章添加到结果列表 e 中
				*e = append(*e, essay)
				break // 找到匹配的关键字后，可以跳出内层循环
			}
		}
	}
}

func GetSearchKeywordRank(rankKind *models.RankKindForZset) (err error) {
	//	得到年月日的keywords的zset
	return redis.GetSearchKeywordRank(rankKind)
}
