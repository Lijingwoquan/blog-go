package models

type Panel struct {
	IpSet    UserIpForSet    `json:"ipSet"`
	RankZset RankKindForZset `json:"rankList"`
}

type UserIpForSet struct {
	Year  int64 `json:"year"`
	Month int64 `json:"month"`
	Week  int64 `json:"week"`
}

type RankKindForZset struct {
	Year  RankListForZset `json:"year"`
	Month RankListForZset `json:"month"`
	Week  RankListForZset `json:"week"`
}

type RankListForZset struct {
	X []string `json:"x"`
	Y []int    `json:"y"`
}
