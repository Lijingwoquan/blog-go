package models

type IndexData struct {
	KindList  []KindData  `json:"kindList"`
	LabelList []LabelData `json:"labelList"`
	EssayList []EssayData `json:"essayList"`
}
type KindData struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Icon       string `json:"icon" db:"icon"`
	EssayCount int8   `json:"essayCount" db:"essay_count"`
}
type LabelData struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name"  db:"name"`
	Color string `json:"color"`
}
