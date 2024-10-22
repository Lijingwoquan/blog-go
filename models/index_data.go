package models

type IndexData struct {
	KindList  []KindData  `json:"kindList"`
	LabelList []LabelData `json:"labelList"`
	EssayList []EssayData `json:"essayList"`
}
