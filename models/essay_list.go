package models

type EssayData struct {
	ID           int         `json:"id" db:"id"`
	Name         string      `json:"name" db:"name"`
	LabelList    []LabelData `json:"labelList,omitempty"`
	KindName     string      `json:"kindName,omitempty" db:"kind_name"`
	KindID       int         `json:"kindID,omitempty" db:"kind_id"`
	Introduction string      `json:"introduction,omitempty" db:"introduction"`
	CreatedTime  string      `json:"createdTime" db:"created_time"`
	VisitedTimes int64       `json:"visitedTimes,omitempty" db:"visited_times"`
	Content      string      `json:"content,omitempty" db:"content"`
	ImgUrl       string      `json:"imgUrl" db:"img_url"`
	Keywords     []string    `json:"keywords,omitempty"`
	IfRecommend  bool        `json:"ifRecommend" db:"if_recommend"`
	IfTop        bool        `json:"ifTop" db:"if_top"`
}

// EssayListAndPage 分页查询相关
type EssayListAndPage struct {
	EssayList []EssayData `json:"essayList,omitempty"`
	TotalPage int         `json:"totalPage,omitempty"`
}
