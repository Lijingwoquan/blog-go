package models

type DataAboutIndex struct {
	ClassifyKindName string                     `json:"classifyKind"`
	ClassifyDetails  []DataAboutClassifyDetails `json:"classifyDetails"`
}

type DataAboutClassifyDetails struct {
	ClassifyKindName string `json:"classifyKind" db:"classifyKindName"`
	ClassifyName     string `json:"classifyName"  db:"classifyName"`
	ClassifyRoute    string `json:"classifyRoute" db:"classifyRoute"`
}
