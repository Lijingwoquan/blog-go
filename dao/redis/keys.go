package redis

const (
	Prefix                = "blog:" //项目key前缀
	KeyVisitedTimes       = "post:visited"
	KeyChangeVisitedTimes = "post:changeVisited"
)

func getRedisKey(key string) string {
	return Prefix + key
}