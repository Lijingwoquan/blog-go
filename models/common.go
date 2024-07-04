package models

type RankKind struct {
	Year  []RankList `json:"year"`
	Month []RankList `json:"month"`
	Week  []RankList `json:"week"`
}

type RankList struct {
	Keyword string `json:"keyword"`
	Times   int    `json:"times"`
}
