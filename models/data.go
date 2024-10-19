package models

// IndexData 返回数据相关
type IndexData struct {
	KindList  []KindData       `json:"kindList"`
	LabelList []LabelData      `json:"labelList"`
	EssayList []DataAboutEssay `json:"essayList"`
}
type KindData struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Icon       string `json:"icon" db:"icon"`
	EssayCount int8   `json:"essayCount" db:"essayCount"`
}
type LabelData struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name"  db:"name"`
	Color string `json:"color"`
}
type DataAboutEssay struct {
	ID           int         `json:"id" db:"id"`
	Name         string      `json:"name" db:"name"`
	LabelList    []LabelData `json:"label_list,omitempty"`
	KindName     string      `json:"kind_name,omitempty" db:"kind_name"`
	KindID       int         `json:"kind_id" db:"kind_id"`
	Introduction string      `json:"introduction,omitempty" db:"introduction"`
	CreatedTime  string      `json:"createdTime" db:"createdTime"`
	VisitedTimes int64       `json:"visitedTimes,omitempty" db:"visitedTimes"`
	Content      string      `json:"content,omitempty" db:"content"`
	ImgUrl       string      `json:"imgUrl" db:"imgUrl"`
	Keywords     []string    `json:"keywords,omitempty"`
	IfRecommend  bool        `json:"ifRecommend" db:"ifRecommend"`
}

// DataAboutEssayListAndPage 分页查询相关
type DataAboutEssayListAndPage struct {
	EssayList []DataAboutEssay `json:"essay_list,omitempty"`
	TotalPage int              `json:"totalPage,omitempty"`
}

// EssayData文章查询相关

type EssayData struct {
	Name         string        `json:"name" db:"name"`
	Kind         string        `json:"kind" db:"kind"`
	Id           int           `json:"id" db:"id"`
	Introduction string        `json:"introduction" db:"introduction"`
	Content      string        `json:"content" db:"content"`
	VisitedTimes int64         `json:"visitedTimes" db:"visitedTimes"`
	CreatedTime  string        `json:"createdTime" db:"createdTime"`
	Keywords     []string      `json:"keywords"`
	Eid          int64         `json:"eid" db:"eid"`
	ImgUrl       string        `json:"imgUrl" binging:"required"  db:"imgUrl"`
	Next         AdjacentEssay `json:"next"`
	Last         AdjacentEssay `json:"last"`
}

type AdjacentEssay struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
