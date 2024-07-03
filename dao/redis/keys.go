package redis

const (
	Prefix                = "blog:" //项目key前缀
	KeyVisitedTimes       = "post:visited:"
	KeyChangeVisitedTimes = "post:changeVisited:"
	KeySearchKeyWordTimes = "keyword:searchTimes:"
	KeyEssayKeyword       = "essay:keyword:"
	KeyUserIp             = "user:ip:"
)

func getRedisKey(key string) string {
	return Prefix + key
}
