package models

type SearchParam struct {
	Keyword string `json:"keyword" binging:"required"`
	IfAdd   bool   `json:"ifAdd"`
}

type EssayIdAndKeyword struct {
	EssayId  int      `json:"essayId" binging:"required"`
	Keywords []string `json:"keywords"`
}
