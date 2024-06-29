package models

type Keyword struct {
	Keyword string `json:"keyword"`
}

type EssayKeyword struct {
	EssayId  int      `json:"essayId"`
	Keywords []string `json:"keywords"`
}
