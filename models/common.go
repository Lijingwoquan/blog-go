package models

type RankKindForZset struct {
	Year  RankListForZset `json:"year"`
	Month RankListForZset `json:"month"`
	Week  RankListForZset `json:"week"`
}

type RankListForZset struct {
	X []string `json:"x"`
	Y []int    `json:"y"`
}

type SetKind struct {
	Year  int64 `json:"year"`
	Month int64 `json:"month"`
	Week  int64 `json:"week"`
}
