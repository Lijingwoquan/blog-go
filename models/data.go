package models

type DataAboutIndex struct {
	DataAboutIndexMenu []DataAboutIndexMenu `json:"dataAboutIndexMenu"`
}

type DataAboutIndexMenu struct {
	DataAboutKind
	Classify []DataAboutClassify `json:"classifyDetails"`
}

type DataAboutKind struct {
	ClassifyKind string `json:"classifyKind" db:"name"`
	Icon         string `json:"icon" db:"icon"`
	Id           int    `json:"id" db:"id"`
}

type DataAboutClassify struct {
	Kind   string           `json:"kind" db:"kind"`
	Name   string           `json:"name"  db:"name"`
	Router string           `json:"router" db:"router"`
	ID     int              `json:"id" db:"id"`
	Essay  []DataAboutEssay `json:"essay" db:"name"`
}

type DataAboutEssay struct {
	Name         string `json:"name" db:"name"`
	Kind         string `json:"kind" db:"kind"`
	Router       string `json:"router" db:"router"`
	Introduction string `json:"introduction" db:"introduction"`
	ID           int    `json:"id" db:"id"`
	CreatedTime  string `json:"createdTime" db:"createdTime"`
	Page         int    `json:"page"` //返回文章对应的页面 实现分页操作
}

type EssayData struct {
	Name         string `json:"name" db:"name"`
	Kind         string `json:"kind" db:"kind"`
	Id           int    `json:"id" db:"id"`
	Introduction string `json:"introduction" db:"introduction"`
	Content      string `json:"content" db:"content"`
	VisitedTimes int64  `json:"visitedTimes" db:"visitedTimes"`
	CreatedTime  string `json:"createdTime" db:"createdTime"`
	UpdatedTime  string `json:"updatedTime" db:"updatedTime"`
}
