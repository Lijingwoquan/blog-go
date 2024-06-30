package models

type Keyword struct {
	Keyword string `json:"keyword"`
}

type EssayKeyword struct {
	EssayId  int      `json:"essayId" binging:"required"`
	Keywords []string `json:"keywords" binding:"omitempty"`
}
