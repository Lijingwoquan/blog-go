package models

// DataAboutIndex 返回数据相关
type DataAboutIndex struct {
	KindList  []DataAboutKind  `json:"kindList"`
	LabelList []DataAboutLabel `json:"labelList"`
	EssayList []DataAboutEssay `json:"essayList"`
}

type DataAboutKind struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Icon       string `json:"icon" db:"icon"`
	EssayCount int8   `json:"essayCount" db:"essayCount"`
}

type DataAboutLabel struct {
	ID   int    `json:"id" db:"id"`
	Kind string `json:"kind" db:"kind"`
	Name string `json:"name"  db:"name"`
}

type DataAboutEssay struct {
	ID           int      `json:"id" db:"id"`
	Name         string   `json:"name" db:"name"`
	Kind         string   `json:"kind,omitempty" db:"kind"`
	Label        string   `json:"label" db:"label"`
	Introduction string   `json:"introduction,omitempty" db:"introduction"`
	CreatedTime  string   `json:"createdTime" db:"createdTime"`
	VisitedTimes int64    `json:"visitedTimes,omitempty" db:"visitedTimes"`
	ImgUrl       string   `json:"imgUrl" db:"imgUrl"`
	Eid          int64    `json:"eid,omitempty" db:"eid"`
	Keywords     []string `json:"keywords,omitempty"`
	IfRecommend  bool     `json:"ifRecommend" db:"ifRecommend"`
}
