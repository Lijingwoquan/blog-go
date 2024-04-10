package redis

const (
	Prefix          = "blog:" //项目key前缀
	KeyVisitedTimes = "post:visited"
)

func getRedisKey(key string) string {
	return Prefix + key
}
