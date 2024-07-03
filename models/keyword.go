package models

type KeywordParam struct {
	Keyword string `json:"keyword" binging:"required"`
}

type EssayIdAndKeyword struct {
	EssayId  int      `json:"essayId" binging:"required"`
	Keywords []string `json:"keywords" binding:"omitempty"`
}

type KeywordRankKind struct {
	Year  []KeywordRankList `json:"year"`
	Month []KeywordRankList `json:"month"`
	Week  []KeywordRankList `json:"week"`
}

type KeywordRankList struct {
	Keyword string `json:"keyword"`
	Times   int    `json:"times"`
}
