package models

type KeywordParam struct {
	Keyword string `json:"keyword" binging:"required"`
}

type EssayIdAndKeyword struct {
	EssayId  int      `json:"essayId" binging:"required"`
	Keywords []string `json:"keywords" binding:"omitempty"`
}
